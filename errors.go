package authgoblue

import "errors"

var (
	// Config
	ErrSecretRequired = errors.New("github.com/qwerius/authgoblue: secret is required")
	ErrIssuerRequired = errors.New("github.com/qwerius/authgoblue: issuer is required")

	// JWT
	ErrInvalidToken      = errors.New("github.com/qwerius/authgoblue: invalid token")
	ErrExpiredToken      = errors.New("github.com/qwerius/authgoblue: token has expired")
	ErrInvalidClaims     = errors.New("github.com/qwerius/authgoblue: invalid claims")
	ErrTokenTypeMismatch = errors.New("github.com/qwerius/authgoblue: invalid token type")

	// Authorization
	ErrUnauthorized = errors.New("github.com/qwerius/authgoblue: unauthorized")
	ErrForbidden    = errors.New("github.com/qwerius/authgoblue: forbidden")

	// Context
	ErrUserNotFound = errors.New("github.com/qwerius/authgoblue: user not found in context")

	// Password
	ErrInvalidPassword = errors.New("github.com/qwerius/authgoblue: invalid password")
	ErrPasswordHash    = errors.New("github.com/qwerius/authgoblue: failed to hash password")
)
