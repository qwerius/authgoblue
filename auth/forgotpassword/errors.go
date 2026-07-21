package forgotpassword

import "errors"

var (
	ErrUserNotFound = errors.New(
		"authgoblue: user not found",
	)

	ErrCannotCreateResetToken = errors.New(
		"authgoblue: cannot create reset token",
	)
)
