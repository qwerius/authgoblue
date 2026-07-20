package token

import "errors"

var (
	ErrInvalidIssuer = errors.New(
		"authgoblue: invalid issuer",
	)

	ErrMissingExpiration = errors.New(
		"authgoblue: missing expiration",
	)

	ErrTokenExpired = errors.New(
		"authgoblue: token expired",
	)

	ErrInvalidAccessTokenType = errors.New(
		"authgoblue: invalid access token type",
	)

	ErrInvalidRefreshTokenType = errors.New(
		"authgoblue: invalid refresh token type",
	)

	ErrMissingTokenID = errors.New(
		"authgoblue: missing token id",
	)

	ErrInvalidToken = errors.New(
		"authgoblue: invalid token",
	)
)
