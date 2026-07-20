package token

import (
	"crypto/hmac"
)

func verifyHS256(
	message []byte,
	signature []byte,
	secret []byte,
) bool {

	expected := signHS256(
		message,
		secret,
	)

	return hmac.Equal(
		expected,
		signature,
	)
}
