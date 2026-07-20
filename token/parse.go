package token

import (
	"authgoblue/claims"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	errUnexpectedSigningMethod = errors.New("authgoblue: unexpected signing method")
	errInvalidToken            = errors.New("authgoblue: invalid token")
	errExpiredToken            = errors.New("authgoblue: token expired")
	errInvalidAccessTokenType  = errors.New("authgoblue: invalid access token type")
	errInvalidRefreshTokenType = errors.New("authgoblue: invalid refresh token type")
)

func splitJWT(token string) (string, string, string, bool) {

	first := strings.IndexByte(token, '.')

	if first < 0 {
		return "", "", "", false
	}

	second := strings.IndexByte(
		token[first+1:],
		'.',
	)

	if second < 0 {
		return "", "", "", false
	}

	second += first + 1

	return token[:first],
		token[first+1 : second],
		token[second+1:],
		true
}

func (s *Service) Parse(tokenString string) (claims.Claims, error) {

	var c claims.Claims

	header, payload, signature, ok := splitJWT(tokenString)

	if !ok {
		return c, errInvalidToken
	}

	// header
	// header validation
	if header != jwtHeader {
		return c, errUnexpectedSigningMethod
	}

	// verify signature

	mac := hmac.New(
		sha256.New,
		s.secret,
	)

	mac.Write([]byte(header))
	mac.Write([]byte("."))
	mac.Write([]byte(payload))

	signatureBytes, err := decodeSegment(signature)

	if err != nil {
		return c, errInvalidToken
	}

	if !hmac.Equal(
		mac.Sum(nil),
		signatureBytes,
	) {
		return c, errInvalidToken
	}

	// payload

	payloadBytes, err := decodeSegment(payload)

	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(payloadBytes, &c); err != nil {
		return c, err
	}

	// expiry check

	if c.ExpiresAt > 0 {

		if time.Now().Unix() > c.ExpiresAt {
			return c, errExpiredToken
		}

	}

	return c, nil
}

func (s *Service) parseByType(
	tokenString string,
	tokenType claims.TokenType,
	errInvalidType error,
) (
	claims.Claims,
	error,
) {

	c, err := s.Parse(tokenString)

	if err != nil {
		return c, err
	}

	if c.TokenType != tokenType {
		return c, errInvalidType
	}

	return c, nil
}

func (s *Service) ParseAccessToken(
	tokenString string,
) (
	claims.Claims, error,
) {

	return s.parseByType(
		tokenString,
		claims.TokenTypeAccess,
		errInvalidAccessTokenType,
	)
}

func (s *Service) ParseRefreshToken(
	tokenString string,
) (
	claims.Claims, error,
) {

	return s.parseByType(
		tokenString,
		claims.TokenTypeRefresh,
		errInvalidRefreshTokenType,
	)
}
