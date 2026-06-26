package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	db *sql.DB
}

type Transaction struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}

type ImportRowResult struct {
	Line          int      `json:"line"`
	Valid         bool     `json:"valid"`
	Date          string   `json:"date,omitempty"`
	Label         string   `json:"label,omitempty"`
	Amount        float64  `json:"amount,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	Errors        []string `json:"errors,omitempty"`
	TransactionID string   `json:"transactionId,omitempty"`
	Duplicate     bool     `json:"duplicate,omitempty"`
}

type ImportReport struct {
	ID           string            `json:"id"`
	Status       string            `json:"status"`
	Filename     string            `json:"filename"`
	DetectedRows int              `json:"detectedRows"`
	ValidRows    int              `json:"validRows"`
	InvalidRows  int              `json:"invalidRows"`
	Rows         []ImportRowResult `json:"rows"`
}

func main() {
	port := getenv("API_PORT", "8080")
	db, err := sql.Open("pgx", getenv("DATABASE_URL", "postgres://noversia:noversia@localhost:5432/noversia?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	app := &App{db: db}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", app.healthHandler)
	mux.HandleFunc("GET /api/v1/accounts", app.accountsHandler)
	mux.HandleFunc("GET /api/v1/transactions", app.transactionsHandler)
	mux.HandleFunc("POST /api/v1/transactions/import", app.importTransactionsHandler)
	mux.HandleFunc("GET /api/v1/imports/{id}", app.importReportHandler)
	mux.HandleFunc("GET /api/v1/recommendations", app.recommendationsHandler)
	mux.HandleFunc("POST /api/v1/chat", app.chatHandler)

	log.Printf("Noversia API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "noversia-api", "version": "0.4.0"})
}

func (a *App) accountsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.QueryContext(r.Context(), `SELECT id::text, name, type, currency, balance::float8 FROM accounts ORDER BY created_at`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	defer rows.Close()

	type Account struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Currency string `json:"currency"`
		Balance float64 `json:"balance"`
	}

	accounts := []Account{}
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Currency, &acc.Balance); err != nil {
			writeError(w, http.StatusInternalServerError, "DB_SCAN_ERROR", err.Error())
			return
		}
		accounts = append(accounts, acc)
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (a *App) transactionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.QueryContext(r.Context(), `
		SELECT id::text, label, amount::float8, currency, booked_at::text
		FROM transactions
		ORDER BY booked_at DESC, created_at DESC
		LIMIT 200`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.ID, &tx.Label, &tx.Amount, &tx.Currency, &tx.Date); err != nil {
			writeError(w, http.StatusInternalServerError, "DB_SCAN_ERROR", err.Error())
			return
		}
		transactions = append(transactions, tx)
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (a *App) importTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_MULTIPART", "Formulaire multipart invalide")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "MISSING_FILE", "Le fichier CSV est obligatoire")
		return
	}
	defer file.Close()

	report, err := parseTransactionCSV(file, header.Filename)
	if err != nil {
		writeError(w, http.StatusBadRequest, "CSV_PARSE_ERROR", err.Error())
		return
	}

	if err := a.persistImportReport(r.Context(), &report); err != nil {
		writeError(w, http.StatusInternalServerError, "IMPORT_PERSIST_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, report)
}

func parseTransactionCSV(reader io.Reader, filename string) (ImportReport, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return ImportReport{}, err
	}

	report := ImportReport{Status: "validated", Filename: filename, Rows: []ImportRowResult{}}
	if len(records) == 0 {
		return report, nil
	}

	header := map[string]int{}
	for i, col := range records[0] {
		header[strings.ToLower(strings.TrimSpace(col))] = i
	}

	for _, col := range []string{"date", "label", "amount", "currency"} {
		if _, ok := header[col]; !ok {
			return ImportReport{}, fmt.Errorf("colonne obligatoire manquante: %s", col)
		}
	}

	for idx, record := range records[1:] {
		line := idx + 2
		result := ImportRowResult{Line: line, Valid: true}

		dateValue := getCSVValue(record, header["date"])
		labelValue := getCSVValue(record, header["label"])
		amountValue := getCSVValue(record, header["amount"])
		currencyValue := strings.ToUpper(getCSVValue(record, header["currency"]))

		if dateValue == "" {
			result.Errors = append(result.Errors, "date obligatoire")
		} else if _, err := time.Parse("2006-01-02", dateValue); err != nil {
			result.Errors = append(result.Errors, "date invalide, format attendu YYYY-MM-DD")
		} else {
			result.Date = dateValue
		}

		if labelValue == "" {
			result.Errors = append(result.Errors, "label obligatoire")
		} else {
			result.Label = labelValue
		}

		amount, err := strconv.ParseFloat(strings.ReplaceAll(amountValue, ",", "."), 64)
		if amountValue == "" {
			result.Errors = append(result.Errors, "amount obligatoire")
		} else if err != nil {
			result.Errors = append(result.Errors, "amount invalide")
		} else {
			result.Amount = amount
		}

		if currencyValue == "" {
			result.Errors = append(result.Errors, "currency obligatoire")
		} else {
			result.Currency = currencyValue
		}

		if len(result.Errors) > 0 {
			result.Valid = false
			report.InvalidRows++
		} else {
			report.ValidRows++
		}

		report.Rows = append(report.Rows, result)
	}

	report.DetectedRows = len(records) - 1
	return report, nil
}

func (a *App) persistImportReport(ctx context.Context, report *ImportReport) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID string
	if err := tx.QueryRowContext(ctx, `SELECT id::text FROM users WHERE email = 'demo@noversia.com'`).Scan(&userID); err != nil {
		return err
	}

	var accountID string
	if err := tx.QueryRowContext(ctx, `SELECT id::text FROM accounts WHERE user_id = $1 ORDER BY created_at LIMIT 1`, userID).Scan(&accountID); err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, `
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

		var transactionID sql.NullString
		if row.Valid {
			sourceHash := hashTransaction(row.Date, row.Label, row.Amount, row.Currency)
			err := tx.QueryRowContext(ctx, `
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

		_, err = tx.ExecContext(ctx, `
			INSERT INTO import_rows (import_batch_id, line_number, valid, raw_data, errors, transaction_id)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::uuid)`,
			report.ID, row.Line, row.Valid, rawData, errorsJSON, row.TransactionID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (a *App) importReportHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var report ImportReport
	err := a.db.QueryRowContext(r.Context(), `
		SELECT id::text, status, filename, detected_rows, valid_rows, invalid_rows
		FROM import_batches WHERE id = $1`, id,
	).Scan(&report.ID, &report.Status, &report.Filename, &report.DetectedRows, &report.ValidRows, &report.InvalidRows)
	if err != nil {
		writeError(w, http.StatusNotFound, "IMPORT_NOT_FOUND", "Import introuvable")
		return
	}

	rows, err := a.db.QueryContext(r.Context(), `
		SELECT line_number, valid, raw_data, errors, COALESCE(transaction_id::text, '')
		FROM import_rows WHERE import_batch_id = $1 ORDER BY line_number`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	defer rows.Close()

	report.Rows = []ImportRowResult{}
	for rows.Next() {
		var row ImportRowResult
		var rawData []byte
		var errorsJSON []byte
		var txID string
		if err := rows.Scan(&row.Line, &row.Valid, &rawData, &errorsJSON, &txID); err != nil {
			writeError(w, http.StatusInternalServerError, "DB_SCAN_ERROR", err.Error())
			return
		}
		_ = json.Unmarshal(rawData, &row)
		_ = json.Unmarshal(errorsJSON, &row.Errors)
		row.TransactionID = txID
		report.Rows = append(report.Rows, row)
	}
	writeJSON(w, http.StatusOK, report)
}

func (a *App) recommendationsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []map[string]any{
		{"id": "rec_001", "title": "Vérifier les abonnements", "description": "Un abonnement récurrent a été détecté.", "confidenceScore": 0.82},
	})
}

func (a *App) chatHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"answer": "Analyse IA simulée : vos dépenses principales semblent concentrées sur courses, abonnements et dépenses variables.",
		"confidenceScore": 0.64,
		"source": "mock",
	})
}

func getCSVValue(record []string, index int) string {
	if index < 0 || index >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[index])
}

func hashTransaction(date string, label string, amount float64, currency string) string {
	payload := fmt.Sprintf("%s|%s|%.2f|%s", date, strings.ToUpper(strings.TrimSpace(label)), amount, strings.ToUpper(currency))
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
