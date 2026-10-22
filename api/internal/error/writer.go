package apperror

import (
	"errors"
	"net/http"

	"github.com/pkab00/dockloud/shared/contract"
)

func WriteBadRequestError(w http.ResponseWriter) {
	contract.WriteBadRequestErrorAsJSON(w)
}

func WriteUnknownError(w http.ResponseWriter) {
	contract.WriteErrorAsJSON(
		w,
		http.StatusInternalServerError,
		"API_ERROR",
		"Unknown API error.",
	)
}

func WriteClientError(w http.ResponseWriter, err ClientError) {
	switch err.Code {
	case contract.ErrInvalidUsername:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Invalid username. Please make sure you provided a correct string.",
		)
	case contract.ErrInvalidEmail:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Invalid email. Please check the format amd try again.",
		)
	case contract.ErrUsernameTaken:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Username is already in use. Please change it and try again.",
		)
	case contract.ErrEmailTaken:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Email is already in use. If you already have an account, try using '/register'.",
		)
	case contract.ErrRegistrationFailed:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Registration service failed with an error.",
		)
	case contract.ErrAuthorizationFailed:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Failed to log in. Please check the data you provided and try again.",
		)
	case contract.ErrUserNotFound:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"No existing users found by request.",
		)
	case contract.ErrUserSearchFailed:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Failed to find user.",
		)
	case contract.ErrDocumentNotFound:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"No existing documents found by request.",
		)
	case contract.ErrInvalidDocumentID:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Invalid document ID passed. Please fix your request and try again.",
		)
	case contract.ErrInvalidDocumentName:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Unsupported document name. Please fix your request and try again.",
		)
	case contract.ErrPermissionDenied:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Permission denied.",
		)
	case contract.ErrLinksNotFound:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Links not found.",
		)
	case contract.ErrLinkExpired:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Link is expired.",
		)
	case contract.ErrLinkCreationNotAllowed:
		contract.WriteErrorAsJSON(
			w,
			http.StatusNotFound,
			contract.ErrDocumentNotFound,
			"No existing documents found by request.",
		)
	default:
		contract.WriteErrorAsJSON(
			w,
			http.StatusInternalServerError,
			err.Code,
			"Unknown client error.",
		)
	}
}

func HandleError(w http.ResponseWriter, err error) {
	if clientErr, ok := errors.AsType[*ClientError](err); ok {
		WriteClientError(w, *clientErr)
	} else {
		WriteUnknownError(w)
	}
}
