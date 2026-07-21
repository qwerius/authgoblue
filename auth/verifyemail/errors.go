package verifyemail

import "errors"

var (
	ErrInvalidVerifyToken = errors.New(
		"authgoblue: invalid verify token",
	)
)
