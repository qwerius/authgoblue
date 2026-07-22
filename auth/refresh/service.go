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
		token:   tokenService,
		revoke:  revokeService,
		session: sessionService,
		hooks:   hookRegistry,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	accessToken, refreshToken, refreshExpiresAt, err :=
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

		AccessExpiresAt: claimsData.ExpiresAt,

		RefreshExpiresAt: refreshExpiresAt,

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
	int64,
	error,
) {

	oldClaims, err :=
		s.token.ParseRefreshToken(
			refreshToken,
		)

	if err != nil {
		return "", "", 0, err
	}

	err =
		s.token.ValidateRefreshToken(
			oldClaims,
		)

	if err != nil {
		return "", "", 0, err
	}

	if oldClaims.SessionID == "" {
		return "", "", 0, ErrMissingSessionID
	}

	if oldClaims.TokenID == "" {
		return "", "", 0, ErrMissingTokenID
	}

	if s.revoke != nil {

		revoked, err :=
			s.revoke.IsRevoked(
				oldClaims.TokenID,
			)

		if err != nil {
			return "", "", 0, err
		}

		if revoked {
			return "", "", 0, ErrRefreshTokenReuse
		}
	}

	currentSession, err :=
		s.session.Get(
			oldClaims.SessionID,
		)

	if err != nil {
		return "", "", 0, err
	}

	if currentSession.Revoked {
		return "", "", 0, session.ErrSessionRevoked
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
		return "", "", 0, err
	}

	newRefresh, err :=
		s.token.GenerateRefreshToken(
			newClaims,
		)

	if err != nil {
		return "", "", 0, err
	}

	newRefreshClaims, err :=
		s.token.ParseRefreshToken(
			newRefresh,
		)

	if err != nil {
		return "", "", 0, err
	}

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
			return "", "", 0, err
		}
	}

	err =
		s.session.Touch(
			oldClaims.SessionID,
		)

	if err != nil {
		return "", "", 0, err
	}

	return newAccess, newRefresh, newRefreshClaims.ExpiresAt, nil
}
