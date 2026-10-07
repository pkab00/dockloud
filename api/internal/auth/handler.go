package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var logger = log.Default()

type Handler struct {
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req *ExternalCreateRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	u, err := h.serv.CreateUser(ctx, req)
	if err != nil {
		logger.Printf("error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
