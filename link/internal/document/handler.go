package document

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkab00/dockloud/link/internal/client/data"
	apperror "github.com/pkab00/dockloud/link/internal/error"
)

type Handler struct {
	serv *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{serv: service}
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	var err error

	token := r.PathValue("token")
	if token == "" {
		apperror.WriteBadRequest(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	doc, err := h.serv.Get(ctx, token)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req data.UpdateDocumentRequest
	var err error

	token := r.PathValue("token")
	if token == "" {
		apperror.WriteBadRequest(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	doc, err := h.serv.Update(ctx, token, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	token := r.PathValue("token")
	if token == "" {
		apperror.WriteBadRequest(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	doc, err := h.serv.Delete(ctx, token)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}
