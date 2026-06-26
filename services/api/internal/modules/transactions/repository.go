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

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Transaction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, label, amount::float8, currency, booked_at::text
		FROM transactions
		ORDER BY booked_at DESC, created_at DESC
		LIMIT 200`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.ID, &tx.Label, &tx.Amount, &tx.Currency, &tx.Date); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (r *Repository) PersistImportReport(ctx context.Context, report *ImportReport) error {
	dbtx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer dbtx.Rollback()

	var userID string
	if err := dbtx.QueryRowContext(ctx, `SELECT id::text FROM users WHERE email = 'demo@noversia.com'`).Scan(&userID); err != nil {
		return err
	}

	var accountID string
	if err := dbtx.QueryRowContext(ctx, `SELECT id::text FROM accounts WHERE user_id = $1 ORDER BY created_at LIMIT 1`, userID).Scan(&accountID); err != nil {
		return err
	}

	err = dbtx.QueryRowContext(ctx, `
		INSERT INTO import_batches (user_id, filename, status, detected_rows, valid_rows, invalid_rows)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text`,
		userID, report.Filename, report.Status, report.DetectedRows, report.ValidRows, report.InvalidRows,
	).Scan(&report.ID)
	if err != nil {
		return err
	}

	for i := range report.Rows {
		row := &report.Rows[i]
		rawData, _ := json.Marshal(row)
		errorsJSON, _ := json.Marshal(row.Errors)

		if row.Valid {
			sourceHash := hashTransaction(row.Date, row.Label, row.Amount, row.Currency)
			var transactionID sql.NullString
			err := dbtx.QueryRowContext(ctx, `
				INSERT INTO transactions (account_id, import_batch_id, booked_at, label, raw_label, amount, currency, source_hash)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT (source_hash) DO NOTHING
				RETURNING id::text`,
				accountID, report.ID, row.Date, row.Label, row.Label, row.Amount, row.Currency, sourceHash,
			).Scan(&transactionID)

			if err == sql.ErrNoRows {
				row.Duplicate = true
			} else if err != nil {
				return err
			} else {
				row.TransactionID = transactionID.String
			}
		}

		_, err = dbtx.ExecContext(ctx, `
			INSERT INTO import_rows (import_batch_id, line_number, valid, raw_data, errors, transaction_id)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::uuid)`,
			report.ID, row.Line, row.Valid, rawData, errorsJSON, row.TransactionID,
		)
		if err != nil {
			return err
		}
	}

	return dbtx.Commit()
}

func (r *Repository) GetImportReport(ctx context.Context, id string) (ImportReport, error) {
	var report ImportReport
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, status, filename, detected_rows, valid_rows, invalid_rows
		FROM import_batches WHERE id = $1`, id,
	).Scan(&report.ID, &report.Status, &report.Filename, &report.DetectedRows, &report.ValidRows, &report.InvalidRows)
	if err != nil {
		return ImportReport{}, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT line_number, valid, raw_data, errors, COALESCE(transaction_id::text, '')
		FROM import_rows WHERE import_batch_id = $1 ORDER BY line_number`, id)
	if err != nil {
		return ImportReport{}, err
	}
	defer rows.Close()

	report.Rows = []ImportRowResult{}
	for rows.Next() {
		var row ImportRowResult
		var rawData []byte
		var errorsJSON []byte
		var txID string
		if err := rows.Scan(&row.Line, &row.Valid, &rawData, &errorsJSON, &txID); err != nil {
			return ImportReport{}, err
		}
		_ = json.Unmarshal(rawData, &row)
		_ = json.Unmarshal(errorsJSON, &row.Errors)
		row.TransactionID = txID
		report.Rows = append(report.Rows, row)
	}
	return report, nil
}

func hashTransaction(date string, label string, amount float64, currency string) string {
	payload := fmt.Sprintf("%s|%s|%.2f|%s", date, strings.ToUpper(strings.TrimSpace(label)), amount, strings.ToUpper(currency))
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
