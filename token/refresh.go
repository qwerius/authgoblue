package token

import (
	"github.com/qwerius/authgoblue/claims"
)

func (s *Service) RefreshAccessToken(
	refreshToken string,
) (string, error) {

	c, err := s.ParseRefreshToken(
		refreshToken,
	)

	if err != nil {
		return "", err
	}

	if err := s.ValidateRefreshToken(
		c,
	); err != nil {

		return "", err
	}

	newClaims := claims.Claims{

		UserID: c.UserID,

		SessionID: c.SessionID,

		Username: c.Username,

		Email: c.Email,

		Role: c.Role,

		Permissions: c.Permissions,
	}

	return s.GenerateAccessToken(
		newClaims,
	)
}
