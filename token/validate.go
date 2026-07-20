package token

import (
	"errors"
	"time"

	"github.com/qwerius/authgoblue/claims"
)

var (
	ErrInvalidAccessType  = errors.New("github.com/qwerius/authgoblue: invalid access token type")
	ErrInvalidRefreshType = errors.New("github.com/qwerius/authgoblue: invalid refresh token type")
)

func (s *Service) ValidateAccessToken(
	c claims.Claims,
) error {

	if err := s.validateCommonClaims(c); err != nil {
		return err
	}

	if c.TokenType != claims.TokenTypeAccess {

		return ErrInvalidAccessType
	}

	return nil
}

func (s *Service) ValidateRefreshToken(
	c claims.Claims,
) error {

	if err := s.validateCommonClaims(c); err != nil {
		return err
	}

	if c.TokenType != claims.TokenTypeRefresh {

		return ErrInvalidRefreshType
	}

	if c.TokenID == "" {
		return ErrMissingTokenID
	}

	return nil
}

func (s *Service) validateCommonClaims(
	c claims.Claims,
) error {

	if err := s.validateIssuer(c); err != nil {
		return err
	}

	if err := s.validateExpiration(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) validateIssuer(
	c claims.Claims,
) error {

	if c.Issuer == "" {

		return ErrInvalidIssuer
	}

	if c.Issuer != s.issuer {

		return ErrInvalidIssuer
	}

	return nil
}

func (s *Service) validateExpiration(
	c claims.Claims,
) error {

	if c.ExpiresAt == 0 {

		return ErrMissingExpiration
	}

	now := time.Now().UTC().Unix()

	if now >= c.ExpiresAt {

		return ErrTokenExpired
	}

	return nil
}
