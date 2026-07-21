package resetpassword

import "errors"

var (
	ErrInvalidResetToken = errors.New(
		"authgoblue: invalid reset token",
	)
)
