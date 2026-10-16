package document

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
	"github.com/pkab00/dockloud/shared/contract"
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

	var documents []Response
	var err error
	authorID := r.URL.Query().Get("author_id")
	requesterId := r.URL.Query().Get("requester_id")

	switch {
	case authorID != "" && requesterId != "":
		documents, err = h.serv.GetByAuthorID(ctx, authorID, requesterId)
	case authorID == "" && requesterId != "":
		documents, err = h.serv.GetAccessible(ctx, requesterId)
	default:
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrQueryParamsMissing,
			"Query parameters missing for get request. Please specify author_id, requester_id, or both.",
		)
		return
	}

	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusInternalServerError,
			contract.ErrInternal,
			"Error requesting documents.",
		)
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
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrInvalidDocumentID,
			"Bad request. Passed invalid document ID.",
		)
		return
	}

	requesterId := r.URL.Query().Get("requester_id")
	if requesterId == "" {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrRequesterIDMissing,
			"Requester ID missing. Make sure to pass it using query parameters.",
		)
		return
	}

	document, err := h.serv.GetByID(ctx, documentID, requesterId)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.DocumentNotFound):
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		case errors.Is(err, apperrors.AccessDeclined):
			contract.WriteErrorAsJSON(
				w,
				http.StatusForbidden,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		default:
			logger.Printf("%v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error requesting a document.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad create request. Failed to create a document.",
		)
		return
	}

	document, err := h.serv.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, apperrors.InvalidDocumentName) {
			contract.WriteErrorAsJSON(
				w,
				http.StatusBadRequest,
				contract.ErrInvalidDocumentName,
				"Invalid document name. Make sure you passed a correct string.",
			)
		} else {
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error creating a document.",
			)
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
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad delete request. Failed to delete a document.",
		)
		return
	}

	requesterId := r.URL.Query().Get("requester_id")
	if requesterId == "" {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrRequesterIDMissing,
			"Requester ID missing. Make sure to pass it using query parameters.",
		)
		return
	}

	document, err := h.serv.Delete(ctx, documentID, requesterId)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.DocumentNotFound):
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		case errors.Is(err, apperrors.AccessDeclined):
			contract.WriteErrorAsJSON(
				w,
				http.StatusForbidden,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		default:
			logger.Printf("%v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error deleting a document.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad update request. Failed to update a document.",
		)
		return
	}

	requesterId := r.URL.Query().Get("requester_id")
	if requesterId == "" {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrRequesterIDMissing,
			"Requester ID missing. Make sure to pass it using query parameters.",
		)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad update request. Failed to update a document.",
		)
		return
	}

	document, err := h.serv.Update(ctx, documentID, requesterId, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.DocumentNotFound):
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		case errors.Is(err, apperrors.AccessDeclined):
			contract.WriteErrorAsJSON(
				w,
				http.StatusForbidden,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		default:
			logger.Printf("%v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error updating a document.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) GetPublicly(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrInvalidDocumentID,
			"Bad request. Passed invalid document ID.",
		)
		return
	}

	document, err := h.serv.GetUnsafe(ctx, documentID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.DocumentNotFound):
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		case errors.Is(err, apperrors.AccessDeclined):
			contract.WriteErrorAsJSON(
				w,
				http.StatusForbidden,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		default:
			logger.Printf("%v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error requesting a document.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) DeletePublicly(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad delete request. Failed to delete a document.",
		)
		return
	}

	document, err := h.serv.DeleteUnsafe(ctx, documentID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.DocumentNotFound):
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		case errors.Is(err, apperrors.AccessDeclined):
			contract.WriteErrorAsJSON(
				w,
				http.StatusForbidden,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		default:
			logger.Printf("%v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error deleting a document.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handler) UpdatePublicly(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	documentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad update request. Failed to update a document.",
		)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrBadRequest,
			"Bad update request. Failed to update a document.",
		)
		return
	}

	document, err := h.serv.UpdateUnsafe(ctx, documentID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.DocumentNotFound):
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		case errors.Is(err, apperrors.AccessDeclined):
			contract.WriteErrorAsJSON(
				w,
				http.StatusForbidden,
				contract.ErrDocumentNotFound,
				"No such document found.",
			)
		default:
			logger.Printf("%v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error updating a document.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}
