package login

import "errors"

var (
	ErrInvalidCredentials = errors.New(
		"github.com/qwerius/authgoblue: invalid credentials",
	)

	ErrUserNotFound = errors.New(
		"github.com/qwerius/authgoblue: user not found",
	)

	ErrLoginDisabled = errors.New(
		"github.com/qwerius/authgoblue: login disabled",
	)
)
