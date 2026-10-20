package document

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkab00/dockloud/api/internal/client/data"
	apperror "github.com/pkab00/dockloud/api/internal/error"
	"github.com/pkab00/dockloud/api/internal/model"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req *ExternalCreateDocumentRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteBadRequestError(w)
		return
	}

	d, err := h.serv.Create(ctx, *req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var d *model.Document
	var req *data.UpdateDocumentRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteBadRequestError(w)
		return
	}

	docID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		apperror.WriteBadRequestError(w)
		return
	}

	d, err = h.serv.Update(ctx, docID, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) UpdateShared(w http.ResponseWriter, r *http.Request) {
	var d *model.Document
	var req *data.UpdateDocumentRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteBadRequestError(w)
		return
	}

	linkToken := r.PathValue("token")

	d, err = h.serv.UpdateShared(ctx, linkToken, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	docID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		apperror.WriteBadRequestError(w)
		return
	}

	d, err := h.serv.Delete(ctx, docID)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) DeleteShared(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	linkToken := r.PathValue("token")

	d, err := h.serv.DeleteShared(ctx, linkToken)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var d []model.Document
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	if authID := query.Get("author_id"); authID != "" {
		d, err = h.serv.GetByAuthor(ctx, authID)
	} else {
		d, err = h.serv.GetAll(ctx)
	}

	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) GetShared(w http.ResponseWriter, r *http.Request) {
	var d *model.Document
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	linkToken := r.PathValue("token")

	d, err = h.serv.GetShared(ctx, linkToken)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	var d *model.Document
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	docID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		apperror.WriteBadRequestError(w)
		return
	}

	d, err = h.serv.GetByID(ctx, docID)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
