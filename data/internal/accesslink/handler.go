package accesslink

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
	return &Handler{serv: serv}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	var al []AccessLink

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	creatorID := r.URL.Query().Get("creator_id")
	documentID := r.URL.Query().Get("document_id")

	switch {
	case creatorID != "" && documentID != "":
		intDocID, err := strconv.Atoi(documentID)
		if err != nil {
			contract.WriteBadRequestErrorAsJSON(w)
			return
		}
		al, err = h.serv.GetAllForDocument(ctx, creatorID, intDocID)
	case creatorID != "":
		al, err = h.serv.GetAll(ctx, creatorID)
	default:
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrQueryParamsMissing,
			"Missing creator_id and/or document_id query parameters.",
		)
		return
	}

	if err != nil {
		if errors.Is(err, apperrors.NoLinksFound) {
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrLinksNotFound,
				"No links found.",
			)
		} else {
			logger.Printf("error getting access links: %v", err)
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrQueryParamsMissing,
				"Unable to get links.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(al)
}

func (h *Handler) GetByToken(w http.ResponseWriter, r *http.Request) {
	var err error
	var al *AccessLink

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	token := r.PathValue("token")
	if token == "" {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	al, err = h.serv.GetByToken(ctx, token)

	switch {
	case errors.Is(err, apperrors.NoLinksFound):
		contract.WriteErrorAsJSON(
			w,
			http.StatusNotFound,
			contract.ErrLinksNotFound,
			"No links found.",
		)
		return
	case errors.Is(err, apperrors.AccessLinkExpired):
		contract.WriteErrorAsJSON(
			w,
			http.StatusForbidden,
			contract.ErrLinkExpired,
			"Access link expired.",
		)
		return
	case err != nil:
		logger.Printf("Error getting access link: %s", err.Error())
		contract.WriteErrorAsJSON(
			w,
			http.StatusInternalServerError,
			contract.ErrInternal,
			"Unable to get link.",
		)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(al)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	var al *AccessLink
	var req CreateLinkRequest

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	requesterID := r.URL.Query().Get("requester_id")
	if requesterID == "" {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrQueryParamsMissing,
			"Missing requester_id query parameter.",
		)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	al, err = h.serv.Create(ctx, requesterID, &req)
	if err != nil {
		logger.Printf("Error creating access link: %v", err)
		contract.WriteErrorAsJSON(
			w,
			http.StatusInternalServerError,
			contract.ErrInternal,
			"Failed to create link.",
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(al)
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	var err error
	var al *AccessLink

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	token := r.PathValue("token")
	if token == "" {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	requesterID := r.URL.Query().Get("requester_id")
	if requesterID == "" {
		contract.WriteErrorAsJSON(
			w,
			http.StatusBadRequest,
			contract.ErrQueryParamsMissing,
			"Missing requester_id query parameter.",
		)
		return
	}

	al, err = h.serv.Revoke(ctx, requesterID, token)
	if err != nil {
		logger.Printf("Failed to revoke link %s.", token)
		contract.WriteErrorAsJSON(
			w,
			http.StatusInternalServerError,
			contract.ErrInternal,
			"Failed to revoke link.",
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(al)
}
