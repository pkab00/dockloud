package apperror

import "errors"

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenRevoked = errors.New("token revoked")
	ErrNoPermission = errors.New("no permission")
)
