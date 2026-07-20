package refresh

import "errors"

var (
	ErrInvalidRefreshToken = errors.New(
		"github.com/qwerius/authgoblue: invalid refresh token",
	)

	ErrRefreshTokenRevoked = errors.New(
		"github.com/qwerius/authgoblue: refresh token revoked",
	)

	ErrRefreshTokenReuse = errors.New(
		"github.com/qwerius/authgoblue: refresh token reuse detected",
	)

	ErrSessionServiceUnavailable = errors.New(
		"github.com/qwerius/authgoblue: session service unavailable",
	)

	ErrMissingSessionID = errors.New(
		"github.com/qwerius/authgoblue: missing session id",
	)
)
