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

func main() {
	port := getenv("API_PORT", "8080")
	db, err := sql.Open("pgx", getenv("DATABASE_URL", "postgres://noversia:noversia@localhost:5432/noversia?sslmode=disable"))
	if err != nil { log.Fatal(err) }
	if err := db.Ping(); err != nil { log.Fatal(err) }

	handler := transactions.NewHandler(transactions.NewService(transactions.NewRepository(db)))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, map[string]any{"status":"ok","version":"0.8.0"}) })
	mux.HandleFunc("GET /api/v1/categories", handler.ListCategories)
	mux.HandleFunc("GET /api/v1/categorization-rules", handler.ListRules)
	mux.HandleFunc("POST /api/v1/categorization-rules", handler.CreateRule)
	mux.HandleFunc("GET /api/v1/rule-suggestions", handler.ListRuleSuggestions)
	mux.HandleFunc("POST /api/v1/rule-suggestions/apply", handler.ApplyRuleSuggestion)
	mux.HandleFunc("GET /api/v1/transactions", handler.List)
	mux.HandleFunc("PATCH /api/v1/transactions/{id}/category", handler.UpdateCategory)
	mux.HandleFunc("POST /api/v1/transactions/import", handler.ImportCSV)
	mux.HandleFunc("GET /api/v1/imports/{id}", handler.GetImportReport)
	log.Printf("Noversia API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil { log.Fatal(err) }
}

func writeJSON(w http.ResponseWriter, status int, payload any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(payload) }
func withCORS(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Access-Control-Allow-Origin", "*"); w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization"); w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS"); if r.Method == http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }; next.ServeHTTP(w, r) }) }
func getenv(key string, fallback string) string { if value := os.Getenv(key); value != "" { return value }; return fallback }
