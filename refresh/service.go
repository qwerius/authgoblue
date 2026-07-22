package refresh

import (
	"time"

	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	token *token.Service

	revoke *revoke.Service

	session *session.Service
}

func NewService(
	tokenService *token.Service,
	revokeService *revoke.Service,
	sessionService *session.Service,
) *Service {

	return &Service{
		token:   tokenService,
		revoke:  revokeService,
		session: sessionService,
	}
}

func (s *Service) Rotate(
	refreshToken string,
) (
	string,
	string,
	claims.Claims,
	error,
) {

	oldClaims, err :=
		s.token.ParseRefreshToken(
			refreshToken,
		)

	if err != nil {
		return "", "", claims.Claims{}, err
	}

	err =
		s.token.ValidateRefreshToken(
			oldClaims,
		)

	if err != nil {
		return "", "", claims.Claims{}, err
	}

	if oldClaims.SessionID == "" {

		return "", "", claims.Claims{}, ErrMissingSessionID
	}

	// cek session masih aktif
	currentSession, err :=
		s.session.Get(
			oldClaims.SessionID,
		)

	if err != nil {
		return "", "", claims.Claims{}, err
	}

	if currentSession.Revoked {

		return "", "", claims.Claims{}, session.ErrSessionRevoked
	}

	// atomic consume refresh token
	// hanya satu request yang boleh memakai token ini
	if s.revoke != nil {

		ok, err :=
			s.revoke.Consume(
				oldClaims.TokenID,
				time.Unix(
					oldClaims.ExpiresAt,
					0,
				),
			)

		if err != nil {
			return "", "", claims.Claims{}, err
		}

		if !ok {

			return "", "", claims.Claims{}, ErrRefreshTokenReuse
		}
	}

	// update aktivitas session
	err =
		s.session.Touch(
			oldClaims.SessionID,
		)

	if err != nil {

		return "", "", claims.Claims{}, err
	}

	newClaims := claims.Claims{

		UserID: oldClaims.UserID,

		SessionID: oldClaims.SessionID,

		Username: oldClaims.Username,

		Email: oldClaims.Email,

		Role: oldClaims.Role,

		Permissions: oldClaims.Permissions,
	}

	newAccess, err :=
		s.token.GenerateAccessToken(
			newClaims,
		)

	if err != nil {

		return "", "", claims.Claims{}, err
	}

	newRefresh, err :=
		s.token.GenerateRefreshToken(
			newClaims,
		)

	if err != nil {

		return "", "", claims.Claims{}, err
	}

	return newAccess, newRefresh, newClaims, nil
}
