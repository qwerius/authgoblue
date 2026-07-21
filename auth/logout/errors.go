package logout

import "errors"

var (
	ErrInvalidToken = errors.New(
		"authgoblue: invalid token",
	)

	ErrInvalidRefreshToken = errors.New(
		"authgoblue: invalid refresh token",
	)

	ErrRevokeToken = errors.New(
		"authgoblue: failed revoke token",
	)

	ErrDeleteSession = errors.New(
		"authgoblue: failed delete session",
	)

	ErrHookExecution = errors.New(
		"authgoblue: failed execute logout hook",
	)
)
