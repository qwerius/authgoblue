package password

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")

	ErrInvalidResetToken = errors.New("invalid reset token")

	ErrExpiredResetToken = errors.New("expired reset token")
)
