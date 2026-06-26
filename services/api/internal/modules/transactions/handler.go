package transactions

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "TRANSACTION_LIST_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) ImportCSV(w http.ResponseWriter, r *http.Request) {
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

	report, err := h.service.ImportCSV(r.Context(), file, header.Filename)
	if err != nil {
		writeError(w, http.StatusBadRequest, "IMPORT_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, report)
}

func (h *Handler) GetImportReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	report, err := h.service.GetImportReport(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "IMPORT_NOT_FOUND", "Import introuvable")
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
