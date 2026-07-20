package token

import (
	"crypto/hmac"
	"crypto/sha256"
)

func signHS256(
	message []byte,
	secret []byte,
) []byte {

	mac := hmac.New(
		sha256.New,
		secret,
	)

	_, _ = mac.Write(message)

	return mac.Sum(nil)
}
