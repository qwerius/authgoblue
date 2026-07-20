package session

import "errors"

var (
	ErrSessionNotFound = errors.New(
		"github.com/qwerius/authgoblue: session not found",
	)

	ErrSessionRevoked = errors.New(
		"github.com/qwerius/authgoblue: session revoked",
	)

	ErrSessionExpired = errors.New(
		"github.com/qwerius/authgoblue: session expired",
	)
)
