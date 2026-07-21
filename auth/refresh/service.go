package refresh

import (
	"context"
	"time"

	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	token *token.Service

	revoke *revoke.Service

	session *session.Service

	hooks *hooks.Registry
}

func New(
	tokenService *token.Service,
	revokeService *revoke.Service,
	sessionService *session.Service,
	hookRegistry *hooks.Registry,
) *Service {

	return &Service{

		token: tokenService,

		revoke: revokeService,

		session: sessionService,

		hooks: hookRegistry,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	accessToken, refreshToken, err :=
		s.Rotate(
			req.RefreshToken,
		)

	if err != nil {
		return nil, err
	}

	claimsData, err :=
		s.token.ParseAccessToken(
			accessToken,
		)

	if err != nil {
		return nil, err
	}

	if s.hooks != nil {

		_ = s.hooks.Fire(
			ctx,
			hooks.EventAfterRefresh,
			hooks.Payload{

				UserID: claimsData.UserID,

				SessionID: claimsData.SessionID,

				Token: accessToken,
			},
		)
	}

	return &Response{

		AccessToken: accessToken,

		RefreshToken: refreshToken,

		Claims: claims.Claims{

			UserID: claimsData.UserID,

			SessionID: claimsData.SessionID,

			Username: claimsData.Username,

			Email: claimsData.Email,

			Role: claimsData.Role,

			Permissions: claimsData.Permissions,
		},
	}, nil
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

	// cek session aktif

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

	// update session activity

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
