package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	apperrors "github.com/pkab00/dockloud/data/internal/errors"
)

type Handler struct {
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{serv: serv}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	users, err := h.serv.GetAll(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	var req UserFindRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := h.serv.Find(ctx, &req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserCreateRequest
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := h.serv.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, apperrors.InavalidNickname) {
			http.Error(w, "Invalid Nickname", http.StatusBadRequest)
		} else if errors.Is(err, apperrors.InvalidEmail) {
			http.Error(w, "Invalid Email Format", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
			http.Error(w, "No Such User Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := h.serv.Update(ctx, uid, &req)
	if err != nil {
		if errors.Is(err, apperrors.UserNotFound) {
			http.Error(w, "No Such User Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
