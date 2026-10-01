package apperrors

import "errors"

var (
	UserNotFound = errors.New("user not found")
	InvalidEmail = errors.New("invalid email")
)
