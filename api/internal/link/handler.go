package link

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkab00/dockloud/api/internal/client/link"
	apperror "github.com/pkab00/dockloud/api/internal/error"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var err error

	docID := r.URL.Query().Get("document_id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.service.Get(ctx, docID)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req link.CreateLinkRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apperror.WriteBadRequestError(w)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.service.Create(ctx, &req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	var err error

	token := r.PathValue("token")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err = h.service.Revoke(ctx, token)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
