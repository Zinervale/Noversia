package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type Transaction struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}

type Recommendation struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	ConfidenceScore float64 `json:"confidenceScore"`
}

type ImportRowResult struct {
	Line     int      `json:"line"`
	Valid    bool     `json:"valid"`
	Date     string   `json:"date,omitempty"`
	Label    string   `json:"label,omitempty"`
	Amount   float64  `json:"amount,omitempty"`
	Currency string   `json:"currency,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

type ImportReport struct {
	Status      string            `json:"status"`
	Filename    string            `json:"filename"`
	DetectedRows int              `json:"detectedRows"`
	ValidRows   int               `json:"validRows"`
	InvalidRows int               `json:"invalidRows"`
	Rows        []ImportRowResult `json:"rows"`
}

func main() {
	port := getenv("API_PORT", "8080")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", healthHandler)
	mux.HandleFunc("GET /api/v1/accounts", accountsHandler)
	mux.HandleFunc("POST /api/v1/accounts", createAccountHandler)
	mux.HandleFunc("GET /api/v1/transactions", transactionsHandler)
	mux.HandleFunc("POST /api/v1/transactions/import", importTransactionsHandler)
	mux.HandleFunc("GET /api/v1/recommendations", recommendationsHandler)
	mux.HandleFunc("POST /api/v1/chat", chatHandler)

	log.Printf("Noversia API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "noversia-api", "version": "0.3.0"})
}

func accountsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []Account{
		{ID: "acc_demo_current", Name: "Compte courant", Type: "checking", Currency: "EUR", Balance: 2450.42},
		{ID: "acc_demo_savings", Name: "Livret", Type: "savings", Currency: "EUR", Balance: 8200.00},
	})
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var input Account
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "JSON invalide")
		return
	}
	if input.Name == "" || input.Type == "" {
		writeError(w, http.StatusBadRequest, "MISSING_REQUIRED_FIELD", "name et type sont obligatoires")
		return
	}
	input.ID = "acc_created_demo"
	if input.Currency == "" {
		input.Currency = "EUR"
	}
	writeJSON(w, http.StatusCreated, input)
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []Transaction{
		{ID: "txn_001", Label: "CARREFOUR MARKET", Amount: -82.31, Currency: "EUR", Date: "2026-06-25"},
		{ID: "txn_002", Label: "SALAIRE", Amount: 2450.00, Currency: "EUR", Date: "2026-06-24"},
		{ID: "txn_003", Label: "NETFLIX", Amount: -13.49, Currency: "EUR", Date: "2026-06-23"},
	})
}

func importTransactionsHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusAccepted, report)
}

func parseTransactionCSV(reader io.Reader, filename string) (ImportReport, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return ImportReport{}, err
	}

	report := ImportReport{
		Status:   "validated",
		Filename: filename,
		Rows:     []ImportRowResult{},
	}

	if len(records) == 0 {
		return report, nil
	}

	header := map[string]int{}
	for i, col := range records[0] {
		header[strings.ToLower(strings.TrimSpace(col))] = i
	}

	required := []string{"date", "label", "amount", "currency"}
	for _, col := range required {
		if _, ok := header[col]; !ok {
			return ImportReport{}, &csvValidationError{message: "colonne obligatoire manquante: " + col}
		}
	}

	for idx, record := range records[1:] {
		line := idx + 2
		result := ImportRowResult{Line: line, Valid: true}

		dateValue := getCSVValue(record, header["date"])
		labelValue := getCSVValue(record, header["label"])
		amountValue := getCSVValue(record, header["amount"])
		currencyValue := getCSVValue(record, header["currency"])

		if strings.TrimSpace(dateValue) == "" {
			result.Errors = append(result.Errors, "date obligatoire")
		} else if _, err := time.Parse("2006-01-02", dateValue); err != nil {
			result.Errors = append(result.Errors, "date invalide, format attendu YYYY-MM-DD")
		} else {
			result.Date = dateValue
		}

		if strings.TrimSpace(labelValue) == "" {
			result.Errors = append(result.Errors, "label obligatoire")
		} else {
			result.Label = labelValue
		}

		amount, err := strconv.ParseFloat(strings.ReplaceAll(amountValue, ",", "."), 64)
		if strings.TrimSpace(amountValue) == "" {
			result.Errors = append(result.Errors, "amount obligatoire")
		} else if err != nil {
			result.Errors = append(result.Errors, "amount invalide")
		} else {
			result.Amount = amount
		}

		if strings.TrimSpace(currencyValue) == "" {
			result.Errors = append(result.Errors, "currency obligatoire")
		} else {
			result.Currency = strings.ToUpper(currencyValue)
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

type csvValidationError struct {
	message string
}

func (e *csvValidationError) Error() string {
	return e.message
}

func getCSVValue(record []string, index int) string {
	if index < 0 || index >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[index])
}

func recommendationsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []Recommendation{
		{
			ID: "rec_001",
			Title: "Vérifier les abonnements",
			Description: "Un abonnement récurrent a été détecté. Il pourra être confirmé ou ignoré dans une prochaine version.",
			ConfidenceScore: 0.82,
		},
	})
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct{ Message string `json:"message"` }
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "JSON invalide")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"answer": "Analyse IA simulée : vos dépenses principales semblent concentrées sur courses, abonnements et dépenses variables.",
		"confidenceScore": 0.64,
		"source": "mock",
	})
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
