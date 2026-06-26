package transactions

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type Repository struct { db *sql.DB }
func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]Transaction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id::text, t.label, t.amount::float8, t.currency, t.booked_at::text,
		       COALESCE(c.id::text, ''), COALESCE(c.name, '')
		FROM transactions t LEFT JOIN categories c ON c.id = t.category_id
		ORDER BY t.booked_at DESC, t.created_at DESC LIMIT 200`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []Transaction{}
	for rows.Next() {
		var item Transaction
		if err := rows.Scan(&item.ID, &item.Label, &item.Amount, &item.Currency, &item.Date, &item.CategoryID, &item.CategoryName); err != nil { return nil, err }
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id::text, name FROM categories ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []Category{}
	for rows.Next() { var item Category; if err := rows.Scan(&item.ID, &item.Name); err != nil { return nil, err }; items = append(items, item) }
	return items, nil
}

func (r *Repository) ListRules(ctx context.Context) ([]CategorizationRule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT r.id::text, r.pattern, r.match_type, r.category_id::text, c.name, r.priority, r.confidence_score::float8, r.enabled
		FROM categorization_rules r JOIN categories c ON c.id = r.category_id
		ORDER BY r.priority ASC, r.created_at ASC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []CategorizationRule{}
	for rows.Next() {
		var item CategorizationRule
		if err := rows.Scan(&item.ID, &item.Pattern, &item.MatchType, &item.CategoryID, &item.CategoryName, &item.Priority, &item.ConfidenceScore, &item.Enabled); err != nil { return nil, err }
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) CreateRule(ctx context.Context, rule CategorizationRule) (CategorizationRule, error) {
	if rule.MatchType == "" { rule.MatchType = "contains" }
	if rule.Priority == 0 { rule.Priority = 100 }
	if rule.ConfidenceScore == 0 { rule.ConfidenceScore = 0.90 }
	rule.Enabled = true
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO categorization_rules (pattern, match_type, category_id, priority, confidence_score, enabled)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id::text`, rule.Pattern, rule.MatchType, rule.CategoryID, rule.Priority, rule.ConfidenceScore).Scan(&rule.ID)
	return rule, err
}

func (r *Repository) UpdateTransactionCategory(ctx context.Context, transactionID string, categoryID string, reason string) (Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil { return Transaction{}, err }
	defer tx.Rollback()

	var previous sql.NullString
	err = tx.QueryRowContext(ctx, `SELECT COALESCE(category_id::text, '') FROM transactions WHERE id = $1`, transactionID).Scan(&previous)
	if err != nil { return Transaction{}, err }

	_, err = tx.ExecContext(ctx, `UPDATE transactions SET category_id = $1, confidence_score = 1.00 WHERE id = $2`, categoryID, transactionID)
	if err != nil { return Transaction{}, err }

	_, err = tx.ExecContext(ctx, `
		INSERT INTO transaction_enrichments (transaction_id, enrichment_type, previous_value, new_value, source, reason)
		VALUES ($1, 'category', $2, $3, 'manual', $4)`,
		transactionID, previous.String, categoryID, reason)
	if err != nil { return Transaction{}, err }

	if err := tx.Commit(); err != nil { return Transaction{}, err }

	items, err := r.List(ctx)
	if err != nil { return Transaction{}, err }
	for _, item := range items {
		if item.ID == transactionID { return item, nil }
	}
	return Transaction{}, sql.ErrNoRows
}

func (r *Repository) ListRuleSuggestions(ctx context.Context) ([]RuleSuggestion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT UPPER(SPLIT_PART(t.label, ' ', 1)) AS pattern,
		       c.id::text,
		       c.name,
		       COUNT(*) AS occurrences
		FROM transaction_enrichments e
		JOIN transactions t ON t.id = e.transaction_id
		JOIN categories c ON c.id::text = e.new_value
		WHERE e.enrichment_type = 'category'
		GROUP BY pattern, c.id, c.name
		HAVING COUNT(*) >= 1
		ORDER BY occurrences DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []RuleSuggestion{}
	for rows.Next() {
		var item RuleSuggestion
		if err := rows.Scan(&item.Pattern, &item.CategoryID, &item.CategoryName, &item.Occurrences); err != nil { return nil, err }
		item.Reason = "Correction manuelle récurrente détectée"
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) PersistImportReport(ctx context.Context, report *ImportReport) error {
	dbtx, err := r.db.BeginTx(ctx, nil)
	if err != nil { return err }
	defer dbtx.Rollback()

	var userID string
	if err := dbtx.QueryRowContext(ctx, `SELECT id::text FROM users WHERE email = 'demo@noversia.com'`).Scan(&userID); err != nil { return err }
	var accountID string
	if err := dbtx.QueryRowContext(ctx, `SELECT id::text FROM accounts WHERE user_id = $1 ORDER BY created_at LIMIT 1`, userID).Scan(&accountID); err != nil { return err }

	err = dbtx.QueryRowContext(ctx, `INSERT INTO import_batches (user_id, filename, status, detected_rows, valid_rows, invalid_rows) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id::text`,
		userID, report.Filename, report.Status, report.DetectedRows, report.ValidRows, report.InvalidRows).Scan(&report.ID)
	if err != nil { return err }

	for i := range report.Rows {
		row := &report.Rows[i]
		if row.Valid {
			sourceHash := hashTransaction(row.Date, row.Label, row.Amount, row.Currency)
			var transactionID sql.NullString
			err := dbtx.QueryRowContext(ctx, `
				INSERT INTO transactions (account_id, category_id, import_batch_id, booked_at, label, raw_label, amount, currency, confidence_score, source_hash)
				VALUES ($1, NULLIF($2, '')::uuid, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (source_hash) DO NOTHING RETURNING id::text`,
				accountID, row.CategoryID, report.ID, row.Date, row.Label, row.Label, row.Amount, row.Currency, row.ConfidenceScore, sourceHash).Scan(&transactionID)
			if err == sql.ErrNoRows { row.Duplicate = true } else if err != nil { return err } else { row.TransactionID = transactionID.String }
		}
		rawData, _ := json.Marshal(row)
		errorsJSON, _ := json.Marshal(row.Errors)
		_, err = dbtx.ExecContext(ctx, `INSERT INTO import_rows (import_batch_id,line_number,valid,raw_data,errors,transaction_id) VALUES ($1,$2,$3,$4,$5,NULLIF($6,'')::uuid)`,
			report.ID, row.Line, row.Valid, rawData, errorsJSON, row.TransactionID)
		if err != nil { return err }
	}
	return dbtx.Commit()
}

func (r *Repository) GetImportReport(ctx context.Context, id string) (ImportReport, error) {
	var report ImportReport
	err := r.db.QueryRowContext(ctx, `SELECT id::text,status,filename,detected_rows,valid_rows,invalid_rows FROM import_batches WHERE id=$1`, id).
		Scan(&report.ID,&report.Status,&report.Filename,&report.DetectedRows,&report.ValidRows,&report.InvalidRows)
	if err != nil { return ImportReport{}, err }
	rows, err := r.db.QueryContext(ctx, `SELECT raw_data FROM import_rows WHERE import_batch_id=$1 ORDER BY line_number`, id)
	if err != nil { return ImportReport{}, err }
	defer rows.Close()
	for rows.Next() { var row ImportRowResult; var raw []byte; if err := rows.Scan(&raw); err != nil { return ImportReport{}, err }; _ = json.Unmarshal(raw,&row); report.Rows = append(report.Rows,row) }
	return report, nil
}

func hashTransaction(date string, label string, amount float64, currency string) string {
	payload := fmt.Sprintf("%s|%s|%.2f|%s", date, strings.ToUpper(strings.TrimSpace(label)), amount, strings.ToUpper(currency))
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
