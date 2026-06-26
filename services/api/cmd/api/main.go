package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/noversia/platform/services/api/internal/modules/transactions"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	db *sql.DB
	transactions *transactions.Handler
}

func main() {
	port := getenv("API_PORT", "8080")
	db, err := sql.Open("pgx", getenv("DATABASE_URL", "postgres://noversia:noversia@localhost:5432/noversia?sslmode=disable"))
	if err != nil { log.Fatal(err) }
	if err := db.Ping(); err != nil { log.Fatal(err) }

	repo := transactions.NewRepository(db)
	service := transactions.NewService(repo)
	handler := transactions.NewHandler(service)

	app := &App{db: db, transactions: handler}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", app.healthHandler)
	mux.HandleFunc("GET /api/v1/accounts", app.accountsHandler)
	mux.HandleFunc("GET /api/v1/categories", app.transactions.ListCategories)
	mux.HandleFunc("GET /api/v1/categorization-rules", app.transactions.ListRules)
	mux.HandleFunc("POST /api/v1/categorization-rules", app.transactions.CreateRule)
	mux.HandleFunc("GET /api/v1/transactions", app.transactions.List)
	mux.HandleFunc("POST /api/v1/transactions/import", app.transactions.ImportCSV)
	mux.HandleFunc("GET /api/v1/imports/{id}", app.transactions.GetImportReport)
	mux.HandleFunc("GET /api/v1/recommendations", app.recommendationsHandler)
	mux.HandleFunc("POST /api/v1/chat", app.chatHandler)

	log.Printf("Noversia API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil { log.Fatal(err) }
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "noversia-api", "version": "0.6.0"})
}

func (a *App) accountsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.QueryContext(r.Context(), `SELECT id::text, name, type, currency, balance::float8 FROM accounts ORDER BY created_at`)
	if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error()); return }
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
			writeError(w, http.StatusInternalServerError, "DB_SCAN_ERROR", err.Error()); return
		}
		accounts = append(accounts, acc)
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (a *App) recommendationsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []map[string]any{
		{"id": "rec_001", "title": "Vérifier les abonnements", "description": "Un abonnement récurrent a été détecté.", "confidenceScore": 0.82},
	})
}

func (a *App) chatHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"answer": "Analyse IA simulée.", "confidenceScore": 0.64, "source": "mock"})
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
		if r.Method == http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }
		next.ServeHTTP(w, r)
	})
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" { return value }
	return fallback
}
