package link

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	apperror "github.com/pkab00/dockloud/link/internal/error"
)

type Handler struct {
	serv *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{serv: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	var req ExternalCreateLinkRequest

	requesterID := r.URL.Query().Get("requester_id")
	if requesterID == "" {
		apperror.WriteBadRequest(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apperror.WriteBadRequest(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	response, err := h.serv.Create(ctx, requesterID, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var err error

	creatorID := r.URL.Query().Get("creator_id")
	documentID := r.URL.Query().Get("document_id")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	response, err := h.serv.GetList(ctx, creatorID, documentID)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	// TODO: auto-revoke link when the document is deleted
	var err error

	requesterID := r.URL.Query().Get("requester_id")
	token := r.PathValue("token")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	response, err := h.serv.Revoke(ctx, token, requesterID)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
