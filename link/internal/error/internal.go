package apperror

import "errors"

var (
	ErrNoRightsToCreate = errors.New("no rights to create link")
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenRevoked     = errors.New("token revoked")
	ErrNoPermission     = errors.New("no permission")
)
