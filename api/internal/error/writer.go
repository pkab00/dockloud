package apperror

import (
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
			"Invalid username. Pleace make sure you provided a correct string.",
		)
	case contract.ErrInvalidEmail:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Invalid email. Pleace check the format amd try again.",
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
	case contract.ErrAccessDeclined:
		contract.WriteErrorAsJSON(
			w,
			err.Status,
			err.Code,
			"Access declined. You have no rights to access this document.",
		)
	default:
		contract.WriteErrorAsJSON(
			w,
			http.StatusInternalServerError,
			err.Code,
			"Unknown internal error.",
		)
	}
}
