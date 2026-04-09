package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponrse struct {
	Error string `json:"error"`
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponrse{Error: msg})
}
