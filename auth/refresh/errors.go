package refresh

import "errors"

var (
	ErrInvalidRefreshToken = errors.New(
		"authgoblue: invalid refresh token",
	)

	ErrRefreshTokenExpired = errors.New(
		"authgoblue: refresh token expired",
	)
)
