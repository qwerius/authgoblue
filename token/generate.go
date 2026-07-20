package token

import (
	"time"

	"github.com/qwerius/authgoblue/claims"

	"github.com/google/uuid"
)

func (s *Service) GenerateAccessToken(
	c claims.Claims,
) (string, error) {

	now := time.Now()

	c.TokenType = claims.TokenTypeAccess

	c.Issuer = s.issuer

	c.IssuedAt = now.Unix()

	c.ExpiresAt = now.Add(
		s.accessTokenTTL,
	).Unix()

	return s.encodeJWT(c)
}

func (s *Service) GenerateRefreshToken(
	c claims.Claims,
) (string, error) {

	now := time.Now()

	c.TokenType = claims.TokenTypeRefresh

	c.TokenID = uuid.NewString()

	c.Issuer = s.issuer

	c.IssuedAt = now.Unix()

	c.ExpiresAt = now.Add(
		s.refreshTokenTTL,
	).Unix()

	return s.encodeJWT(c)
}
