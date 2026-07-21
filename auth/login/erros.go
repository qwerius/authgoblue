package login

import "errors"

var (
	ErrInvalidCredentials = errors.New(
		"authgoblue: invalid credentials",
	)

	ErrAccountDisabled = errors.New(
		"authgoblue: account disabled",
	)

	ErrEmailNotVerified = errors.New(
		"authgoblue: email not verified",
	)
)
