package authgoblue

import "errors"

var (
	// Config
	ErrSecretRequired = errors.New("authgoblue: secret is required")
	ErrIssuerRequired = errors.New("authgoblue: issuer is required")

	// JWT
	ErrInvalidToken      = errors.New("authgoblue: invalid token")
	ErrExpiredToken      = errors.New("authgoblue: token has expired")
	ErrInvalidClaims     = errors.New("authgoblue: invalid claims")
	ErrTokenTypeMismatch = errors.New("authgoblue: invalid token type")

	// Authorization
	ErrUnauthorized = errors.New("authgoblue: unauthorized")
	ErrForbidden    = errors.New("authgoblue: forbidden")

	// Context
	ErrUserNotFound = errors.New("authgoblue: user not found in context")

	// Password
	ErrInvalidPassword = errors.New("authgoblue: invalid password")
	ErrPasswordHash    = errors.New("authgoblue: failed to hash password")
)
