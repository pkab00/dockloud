package document

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

type Handler struct {
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var documents []DocumentResponse
	var err error
	authorID := r.URL.Query().Get("author_id")
	requesterId := r.URL.Query().Get("requester_id")

	switch {
	case authorID != "" && requesterId != "":
		documents, err = h.serv.GetByAuthorID(ctx, authorID, requesterId)
	case authorID == "" && requesterId != "":
		documents, err = h.serv.GetAccessible(ctx, requesterId)
	default:
		http.Error(w, "Query Parameters Missing", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	requesterId := r.URL.Query().Get("requester_id")
	if requesterId == "" {
		http.Error(w, "Requester ID Missing", http.StatusBadRequest)
		return
	}

	document, err := h.serv.GetByID(ctx, documentID, requesterId)
	if err != nil {
		if errors.Is(err, apperrors.DocumentNotFound) {
			http.Error(w, "No Such Document Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req DocumentCreateRequest

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	document, err := h.serv.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, apperrors.InvalidDocumentName) {
			http.Error(w, "Invalid Document Name", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	requesterId := r.URL.Query().Get("requester_id")
	if requesterId == "" {
		http.Error(w, "Requester ID Missing", http.StatusBadRequest)
		return
	}

	document, err := h.serv.Delete(ctx, documentID, requesterId)
	if err != nil {
		if errors.Is(err, apperrors.DocumentNotFound) {
			http.Error(w, "No Such Document Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req DocumentUpdateRequest

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	requesterId := r.URL.Query().Get("requester_id")
	if requesterId == "" {
		http.Error(w, "Requester ID Missing", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	document, err := h.serv.Update(ctx, documentID, requesterId, &req)
	if err != nil {
		if errors.Is(err, apperrors.DocumentNotFound) {
			http.Error(w, "No Such Document Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}
