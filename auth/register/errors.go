package register

import "errors"

var (
	ErrEmailAlreadyExists = errors.New(
		"authgoblue: email already exists",
	)

	ErrUsernameAlreadyExists = errors.New(
		"authgoblue: username already exists",
	)
)
