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

	ErrCreateSession = errors.New(
		"authgoblue: create session failed",
	)

	ErrGenerateToken = errors.New(
		"authgoblue: generate token failed",
	)

	ErrHookExecution = errors.New(
		"authgoblue: hook execution failed",
	)

	ErrInvalidUserID = errors.New(
		"authgoblue: invalid user id",
	)
)
