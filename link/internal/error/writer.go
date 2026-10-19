package apperror

import (
	"errors"
	"net/http"

	"github.com/pkab00/dockloud/shared/contract"
)

func WriteBadRequest(w http.ResponseWriter) {
	contract.WriteBadRequestErrorAsJSON(w)
}

func WriteClientError(w http.ResponseWriter, err ClientError) {
	contract.WriteErrorAsJSON(w, err.Status, err.Code, err.Message)
}

func WriteAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTokenExpired):
		contract.WriteErrorAsJSON(w, http.StatusGone, contract.ErrLinkExpired, "Token expired")
	case errors.Is(err, ErrTokenRevoked):
		contract.WriteErrorAsJSON(w, http.StatusGone, contract.ErrLinkExpired, "Token revoked")
	case errors.Is(err, ErrNoPermission):
		contract.WriteErrorAsJSON(w, http.StatusForbidden, contract.ErrPermissionDenied, "No permission")
	default:
		contract.WriteErrorAsJSON(w, http.StatusInternalServerError, contract.ErrInternal, "Internal server error: "+err.Error())
	}
}

func HandleError(w http.ResponseWriter, err error) {
	if clientErr, ok := errors.AsType[*ClientError](err); ok {
		WriteClientError(w, *clientErr)
	} else {
		WriteAppError(w, err)
	}
}
