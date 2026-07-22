package token

import "errors"

var (
	ErrInvalidIssuer = errors.New(
		"github.com/qwerius/authgoblue: invalid issuer",
	)

	ErrMissingIssuedAt = errors.New(
		"github.com/qwerius/authgoblue: missing issued at",
	)

	ErrInvalidIssuedAt = errors.New(
		"github.com/qwerius/authgoblue: invalid issued at",
	)

	ErrMissingExpiration = errors.New(
		"github.com/qwerius/authgoblue: missing expiration",
	)

	ErrTokenExpired = errors.New(
		"github.com/qwerius/authgoblue: token expired",
	)

	ErrInvalidAccessTokenType = errors.New(
		"github.com/qwerius/authgoblue: invalid access token type",
	)

	ErrInvalidRefreshTokenType = errors.New(
		"github.com/qwerius/authgoblue: invalid refresh token type",
	)

	ErrMissingTokenID = errors.New(
		"github.com/qwerius/authgoblue: missing token id",
	)

	ErrInvalidToken = errors.New(
		"github.com/qwerius/authgoblue: invalid token",
	)
)
