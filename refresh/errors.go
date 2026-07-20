package refresh

import "errors"

var (
	ErrInvalidRefreshToken = errors.New(
		"authgoblue: invalid refresh token",
	)

	ErrRefreshTokenRevoked = errors.New(
		"authgoblue: refresh token revoked",
	)

	ErrRefreshTokenReuse = errors.New(
		"authgoblue: refresh token reuse detected",
	)

	ErrSessionServiceUnavailable = errors.New(
		"authgoblue: session service unavailable",
	)

	ErrMissingSessionID = errors.New(
		"authgoblue: missing session id",
	)
)
