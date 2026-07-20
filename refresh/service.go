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
	error,
) {

	oldClaims, err :=
		s.token.ParseRefreshToken(
			refreshToken,
		)

	if err != nil {
		return "", "", err
	}

	err =
		s.token.ValidateRefreshToken(
			oldClaims,
		)

	if err != nil {
		return "", "", err
	}

	if oldClaims.SessionID == "" {
		return "", "", ErrMissingSessionID
	}

	// cek refresh token reuse
	if s.revoke != nil {

		revoked, err :=
			s.revoke.IsRevoked(
				oldClaims.TokenID,
			)

		if err != nil {
			return "", "", err
		}

		if revoked {
			return "", "", ErrRefreshTokenReuse
		}
	}

	// cek session masih aktif
	currentSession, err :=
		s.session.Get(
			oldClaims.SessionID,
		)

	if err != nil {
		return "", "", err
	}

	if currentSession.Revoked {
		return "", "", session.ErrSessionRevoked
	}

	// revoke refresh token lama
	if s.revoke != nil {

		err =
			s.revoke.Revoke(
				oldClaims.TokenID,
				time.Unix(
					oldClaims.ExpiresAt,
					0,
				),
			)

		if err != nil {
			return "", "", err
		}
	}

	// update aktivitas session
	err =
		s.session.Touch(
			oldClaims.SessionID,
		)

	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	newRefresh, err :=
		s.token.GenerateRefreshToken(
			newClaims,
		)

	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}
