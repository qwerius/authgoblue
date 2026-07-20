package logout

import "errors"

var (
	ErrMissingToken     = errors.New("missing access token")
	ErrInvalidToken     = errors.New("invalid access token")
	ErrMissingSessionID = errors.New("missing session id")
)
