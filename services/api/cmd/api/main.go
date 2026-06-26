package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func main() {
	port := getenv("API_PORT", "8080")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", healthHandler)
	mux.HandleFunc("GET /api/v1/accounts", accountsHandler)
	mux.HandleFunc("POST /api/v1/accounts", createAccountHandler)
	mux.HandleFunc("GET /api/v1/transactions", transactionsHandler)
	mux.HandleFunc("GET /api/v1/recommendations", recommendationsHandler)
	mux.HandleFunc("POST /api/v1/chat", chatHandler)

	log.Printf("Noversia API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "noversia-api",
	})
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
	var input struct {
		Message string `json:"message"`
	}
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
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code": code,
			"message": message,
		},
	})
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
