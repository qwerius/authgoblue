package login

import "errors"

var (
	ErrInvalidCredentials = errors.New(
		"authgoblue: invalid credentials",
	)

	ErrUserNotFound = errors.New(
		"authgoblue: user not found",
	)

	ErrLoginDisabled = errors.New(
		"authgoblue: login disabled",
	)
)
