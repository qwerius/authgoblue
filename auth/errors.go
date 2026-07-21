package auth

import "errors"

var (
	ErrInvalidProvider = errors.New(
		"authgoblue: invalid provider",
	)

	ErrNotConfigured = errors.New(
		"authgoblue: service not configured",
	)

	ErrNotSupported = errors.New(
		"authgoblue: operation not supported",
	)
)
