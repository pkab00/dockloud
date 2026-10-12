package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	var req UserFindRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	res, err := h.serv.Find(ctx, &req)
	if err != nil {
		if errors.Is(err, apperrors.UserNotFound) {
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrUserNotFound,
				"No matching users find by the search request.",
			)
		} else {
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrUserSearchFailed,
				"An error occured while searching for users.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var req UserAuthRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	user, err := h.serv.Authenticate(ctx, &req)
	if err != nil {
		contract.WriteErrorAsJSON(
			w,
			http.StatusUnauthorized,
			contract.ErrAuthorizationFailed,
			"Authorization failed. Please check your log in data.",
		)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserCreateRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	user, err := h.serv.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.InavalidUsername):
			contract.WriteErrorAsJSON(
				w,
				http.StatusBadRequest,
				contract.ErrInvalidUsername,
				"Username is invalid. Please make sure that it contains less than 50 symbols.",
			)
		case errors.Is(err, apperrors.InvalidEmail):
			contract.WriteErrorAsJSON(
				w,
				http.StatusBadRequest,
				contract.ErrInvalidEmail,
				"Invalid or unsupported email format. Please fix your input and try again.",
			)
		case errors.Is(err, apperrors.UsernameAlreadyInUse):
			contract.WriteErrorAsJSON(
				w,
				http.StatusConflict,
				contract.ErrUsernameTaken,
				"Username is already in use. Please change it and try again.",
			)
		case errors.Is(err, apperrors.EmailAlreadyInUse):
			contract.WriteErrorAsJSON(
				w,
				http.StatusConflict,
				contract.ErrEmailTaken,
				"Email is already in use. Please change it and try again.",
			)
		default:
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrRegistrationFailed,
				"Unknown registration error.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	uid := r.PathValue("uid")

	user, err := h.serv.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, apperrors.UserNotFound) {
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrUserNotFound,
				"No users found with such uid",
			)
		} else {
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error trying to delete user.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UserUpdateRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	uid := r.PathValue("uid")
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		contract.WriteBadRequestErrorAsJSON(w)
		return
	}

	user, err := h.serv.Update(ctx, uid, &req)
	if err != nil {
		if errors.Is(err, apperrors.UserNotFound) {
			contract.WriteErrorAsJSON(
				w,
				http.StatusNotFound,
				contract.ErrUserNotFound,
				"No users found with such uid.",
			)
		} else {
			contract.WriteErrorAsJSON(
				w,
				http.StatusInternalServerError,
				contract.ErrInternal,
				"Error trying to delete user.",
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
