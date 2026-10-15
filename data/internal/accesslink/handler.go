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

func NewHandler(serv Service) *Handler {
	return &Handler{serv: &serv}
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
