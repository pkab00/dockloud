package contract

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteErrorAsJSON(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func WriteBadRequestErrorAsJSON(w http.ResponseWriter) {
	WriteErrorAsJSON(
		w,
		http.StatusBadRequest,
		ErrBadRequest,
		"Bad request. Please fix your request's body and try again.",
	)
}

const (
	ErrBadRequest         = "BAD_REQUEST"
	ErrInternal           = "INTERNAL_ERROR"
	ErrReqesterIDMissing  = "REQUESTER_ID_MISSING"
	ErrQueryParamsMissing = "QUERY_PARAMETERS_MISSING"

	ErrUserNotFound        = "USER_NOT_FOUND"
	ErrUserSearchFailed    = "USER_SEARCH_ERROR"
	ErrAuthorizationFailed = "AUTHORIZATION_ERROR"
	ErrInvalidUsername     = "INVALID_USERNAME"
	ErrInvalidEmail        = "INVALID_EMAIL"
	ErrUsernameTaken       = "USERNAME_TAKEN"
	ErrEmailTaken          = "EMAIL_TAKEN"
	ErrRegistrationFailed  = "REGISTRATION_ERROR"

	ErrDocumentNotFound    = "DOCUMENT_NOT_FOUND"
	ErrInvalidDocumentName = "INVALID_DOCUMENT_NAME"
	ErrInvalidDocumentID   = "INVALID_DOCUMENT_ID"

	ErrLinksNotFound = "LINKS_NOT_FOUND"
	ErrLinkExpired   = "LINK_EXPIRED"
)
