package apperror

import "errors"

var (
	UsernameInUse = errors.New("username already in use")
	EmailInUse    = errors.New("email already in use")
)
