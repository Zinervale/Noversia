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
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		ORDER BY t.booked_at DESC, t.created_at DESC
		LIMIT 200`)
	if err != nil { return nil, err }
	defer rows.Close()

	items := []Transaction{}
	for rows.Next() {
		var item Transaction
		if err := rows.Scan(&item.ID, &item.Label, &item.Amount, &item.Currency, &item.Date, &item.CategoryID, &item.CategoryName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id::text, name FROM categories ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var item Category
		if err := rows.Scan(&item.ID, &item.Name); err != nil { return nil, err }
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ListRules(ctx context.Context) ([]CategorizationRule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT r.id::text, r.pattern, r.match_type, r.category_id::text, c.name, r.priority, r.confidence_score::float8, r.enabled
		FROM categorization_rules r
		JOIN categories c ON c.id = r.category_id
		ORDER BY r.priority ASC, r.created_at ASC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := []CategorizationRule{}
	for rows.Next() {
		var item CategorizationRule
		if err := rows.Scan(&item.ID, &item.Pattern, &item.MatchType, &item.CategoryID, &item.CategoryName, &item.Priority, &item.ConfidenceScore, &item.Enabled); err != nil {
			return nil, err
		}
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

func (r *Repository) PersistImportReport(ctx context.Context, report *ImportReport) error {
	dbtx, err := r.db.BeginTx(ctx, nil)
	if err != nil { return err }
	defer dbtx.Rollback()

	var userID string
	if err := dbtx.QueryRowContext(ctx, `SELECT id::text FROM users WHERE email = 'demo@noversia.com'`).Scan(&userID); err != nil { return err }

	var accountID string
	if err := dbtx.QueryRowContext(ctx, `SELECT id::text FROM accounts WHERE user_id = $1 ORDER BY created_at LIMIT 1`, userID).Scan(&accountID); err != nil { return err }

	err = dbtx.QueryRowContext(ctx, `
		INSERT INTO import_batches (user_id, filename, status, detected_rows, valid_rows, invalid_rows)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text`,
		userID, report.Filename, report.Status, report.DetectedRows, report.ValidRows, report.InvalidRows,
	).Scan(&report.ID)
	if err != nil { return err }

	for i := range report.Rows {
		row := &report.Rows[i]
		rawData, _ := json.Marshal(row)
		errorsJSON, _ := json.Marshal(row.Errors)

		if row.Valid {
			sourceHash := hashTransaction(row.Date, row.Label, row.Amount, row.Currency)
			var transactionID sql.NullString
			var categoryID any = nil
			if row.CategoryID != "" { categoryID = row.CategoryID }

			err := dbtx.QueryRowContext(ctx, `
				INSERT INTO transactions (account_id, category_id, import_batch_id, booked_at, label, raw_label, amount, currency, confidence_score, source_hash)
				VALUES ($1, NULLIF($2, '')::uuid, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (source_hash) DO NOTHING
				RETURNING id::text`,
				accountID, categoryID, report.ID, row.Date, row.Label, row.Label, row.Amount, row.Currency, row.ConfidenceScore, sourceHash,
			).Scan(&transactionID)

			if err == sql.ErrNoRows { row.Duplicate = true
			} else if err != nil { return err
			} else { row.TransactionID = transactionID.String }
		}

		rawData, _ = json.Marshal(row)
		_, err = dbtx.ExecContext(ctx, `
			INSERT INTO import_rows (import_batch_id, line_number, valid, raw_data, errors, transaction_id)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::uuid)`,
			report.ID, row.Line, row.Valid, rawData, errorsJSON, row.TransactionID,
		)
		if err != nil { return err }
	}
	return dbtx.Commit()
}

func (r *Repository) GetImportReport(ctx context.Context, id string) (ImportReport, error) {
	var report ImportReport
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, status, filename, detected_rows, valid_rows, invalid_rows
		FROM import_batches WHERE id = $1`, id,
	).Scan(&report.ID, &report.Status, &report.Filename, &report.DetectedRows, &report.ValidRows, &report.InvalidRows)
	if err != nil { return ImportReport{}, err }

	rows, err := r.db.QueryContext(ctx, `SELECT raw_data FROM import_rows WHERE import_batch_id = $1 ORDER BY line_number`, id)
	if err != nil { return ImportReport{}, err }
	defer rows.Close()

	report.Rows = []ImportRowResult{}
	for rows.Next() {
		var row ImportRowResult
		var rawData []byte
		if err := rows.Scan(&rawData); err != nil { return ImportReport{}, err }
		_ = json.Unmarshal(rawData, &row)
		report.Rows = append(report.Rows, row)
	}
	return report, nil
}

func hashTransaction(date string, label string, amount float64, currency string) string {
	payload := fmt.Sprintf("%s|%s|%.2f|%s", date, strings.ToUpper(strings.TrimSpace(label)), amount, strings.ToUpper(currency))
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
