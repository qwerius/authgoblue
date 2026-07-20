package ctx

import (
	"github.com/qwerius/authgoblue/claims"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) TokenType(
	c fiber.Ctx,
) (claims.TokenType, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return "", err
	}

	return authClaims.TokenType, nil
}

func (s *Service) IsAccessToken(
	c fiber.Ctx,
) (bool, error) {

	tokenType, err := s.TokenType(c)

	if err != nil {
		return false, err
	}

	return tokenType == claims.TokenTypeAccess, nil
}

func (s *Service) IsRefreshToken(
	c fiber.Ctx,
) (bool, error) {

	tokenType, err := s.TokenType(c)

	if err != nil {
		return false, err
	}

	return tokenType == claims.TokenTypeRefresh, nil
}

func (s *Service) Issuer(
	c fiber.Ctx,
) (string, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return "", err
	}

	return authClaims.Issuer, nil
}

func (s *Service) IssuedAt(
	c fiber.Ctx,
) (int64, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return 0, err
	}

	return authClaims.IssuedAt, nil
}

func (s *Service) ExpiresAt(
	c fiber.Ctx,
) (int64, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return 0, err
	}

	return authClaims.ExpiresAt, nil
}
