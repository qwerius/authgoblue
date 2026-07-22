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
	int64,
	int64,
	error,
) {

	oldClaims, err :=
		s.token.ParseRefreshToken(
			refreshToken,
		)

	if err != nil {
		return "", "", claims.Claims{}, 0, 0, err
	}

	err =
		s.token.ValidateRefreshToken(
			oldClaims,
		)

	if err != nil {
		return "", "", claims.Claims{}, 0, 0, err
	}

	if oldClaims.SessionID == "" {

		return "", "", claims.Claims{}, 0, 0, ErrMissingSessionID
	}

	// cek session masih aktif
	currentSession, err :=
		s.session.Get(
			oldClaims.SessionID,
		)

	if err != nil {
		return "", "", claims.Claims{}, 0, 0, err
	}

	if currentSession.Revoked {

		return "", "", claims.Claims{}, 0, 0, session.ErrSessionRevoked
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
			return "", "", claims.Claims{}, 0, 0, err
		}

		if !ok {

			return "", "", claims.Claims{}, 0, 0, ErrRefreshTokenReuse
		}
	}

	// update aktivitas session
	err =
		s.session.Touch(
			oldClaims.SessionID,
		)

	if err != nil {

		return "", "", claims.Claims{}, 0, 0, err
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

		return "", "", claims.Claims{}, 0, 0, err
	}

	newRefresh, err :=
		s.token.GenerateRefreshToken(
			newClaims,
		)

	if err != nil {

		return "", "", claims.Claims{}, 0, 0, err
	}

	accessClaims, err :=
		s.token.ParseAccessToken(
			newAccess,
		)

	if err != nil {
		return "", "", claims.Claims{}, 0, 0, err
	}

	refreshClaims, err :=
		s.token.ParseRefreshToken(
			newRefresh,
		)

	if err != nil {
		return "", "", claims.Claims{}, 0, 0, err
	}

	return newAccess,
		newRefresh,
		newClaims,
		accessClaims.ExpiresAt,
		refreshClaims.ExpiresAt,
		nil
}
