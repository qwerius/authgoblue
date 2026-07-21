package refresh

import "errors"

var (
	ErrInvalidRefreshToken = errors.New(
		"authgoblue: invalid refresh token",
	)

	ErrRefreshTokenExpired = errors.New(
		"authgoblue: refresh token expired",
	)

	ErrMissingSessionID = errors.New(
		"authgoblue: missing session id",
	)

	ErrMissingTokenID = errors.New(
		"authgoblue: missing token id",
	)

	ErrRefreshTokenReuse = errors.New(
		"authgoblue: refresh token reuse detected",
	)
)
