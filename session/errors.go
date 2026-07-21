package session

import "errors"

var (
	ErrSessionNotFound = errors.New(
		"authgoblue: session not found",
	)

	ErrSessionRevoked = errors.New(
		"authgoblue: session revoked",
	)

	ErrSessionExpired = errors.New(
		"authgoblue: session expired",
	)
)
