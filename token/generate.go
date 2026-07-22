package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/qwerius/authgoblue/claims"
)

func (s *Service) prepareClaims(
	c claims.Claims,
	tokenType claims.TokenType,
	ttl time.Duration,
) claims.Claims {

	now := time.Now()

	c.TokenType = tokenType
	c.Issuer = s.issuer
	c.IssuedAt = now.Unix()
	c.ExpiresAt = now.Add(ttl).Unix()

	return c
}

func (s *Service) GenerateAccessToken(
	c claims.Claims,
) (string, error) {

	c = s.prepareClaims(
		c,
		claims.TokenTypeAccess,
		s.accessTokenTTL,
	)

	return s.encodeJWT(c)
}

func (s *Service) GenerateRefreshToken(
	c claims.Claims,
) (string, error) {

	c = s.prepareClaims(
		c,
		claims.TokenTypeRefresh,
		s.refreshTokenTTL,
	)

	c.TokenID = uuid.NewString()

	return s.encodeJWT(c)
}
