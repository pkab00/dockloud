package apperrors

import "errors"

var (
	UserNotFound        = errors.New("user not found")
	InavalidNickname    = errors.New("invalid nickname")
	InvalidEmail        = errors.New("invalid email")
	DocumentNotFound    = errors.New("document not found")
	InvalidDocumentName = errors.New("invalid document name")
)
