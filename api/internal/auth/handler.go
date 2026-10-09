package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	apperror "github.com/pkab00/dockloud/api/internal/error"
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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
		switch {
		case errors.Is(err, apperror.UsernameInUse):
			http.Error(w, "Username Already In Use", http.StatusBadRequest)
		case errors.Is(err, apperror.EmailInUse):
			http.Error(w, "Email Already In Use", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req *ExternalAuthRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	res, err := h.serv.AuthUser(ctx, req)
	if err != nil {
		logger.Printf("error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
