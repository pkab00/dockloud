package apperrors

import "errors"

var (
	UserNotFound         = errors.New("user not found")
	AccessDeclined       = errors.New("access declined")
	InavalidUsername     = errors.New("invalid username")
	InvalidEmail         = errors.New("invalid email")
	UsernameAlreadyInUse = errors.New("username is already in use")
	EmailAlreadyInUse    = errors.New("email is already in use")
	DocumentNotFound     = errors.New("document not found")
	InvalidDocumentName  = errors.New("invalid document name")
	NoLinksFound         = errors.New("no links found")
	AccessLinkExpired    = errors.New("access link expired")
)
