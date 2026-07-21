package logout

import "errors"

var (
	ErrInvalidRefreshToken = errors.New(
		"authgoblue: invalid refresh token",
	)
)
