package token

import (
	"authgoblue/claims"
	"encoding/json"
	"errors"
	"strings"
)

var (
	errMalformedToken   = errors.New("authgoblue: malformed token")
	errInvalidSignature = errors.New("authgoblue: invalid signature")
	errInvalidHeader    = errors.New("authgoblue: invalid jwt header")
	errInvalidPayload   = errors.New("authgoblue: invalid payload")
)

func (s *Service) decodeJWT(
	tokenString string,
) (claims.Claims, error) {

	var c claims.Claims

	parts := strings.Split(
		tokenString,
		".",
	)

	if len(parts) != 3 {
		return c, errMalformedToken
	}

	// saat ini hanya mendukung HS256
	if parts[0] != jwtHeader {
		return c, errInvalidHeader
	}

	signature, err := decodeSegment(
		parts[2],
	)

	if err != nil {
		return c, errInvalidSignature
	}

	unsigned := parts[0] + "." + parts[1]

	if !verifyHS256(
		[]byte(unsigned),
		signature,
		s.secret,
	) {
		return c, errInvalidSignature
	}

	payload, err := decodeSegment(
		parts[1],
	)

	if err != nil {
		return c, errInvalidPayload
	}

	if err := json.Unmarshal(
		payload,
		&c,
	); err != nil {

		return c, errInvalidPayload
	}

	return c, nil
}
