package refresh

import (
	"context"

	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	coreRefresh "github.com/qwerius/authgoblue/refresh"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	token *token.Service

	refresh *coreRefresh.Service

	hooks *hooks.Registry
}

func New(
	tokenService *token.Service,
	refreshService *coreRefresh.Service,
	hookRegistry *hooks.Registry,
) *Service {

	return &Service{
		token:   tokenService,
		refresh: refreshService,
		hooks:   hookRegistry,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	accessToken,
		refreshToken,
		tokenClaims,
		err :=
		s.refresh.Rotate(
			req.RefreshToken,
		)

	if err != nil {
		return nil, err
	}

	accessClaims, err :=
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
				UserID: tokenClaims.UserID,

				SessionID: tokenClaims.SessionID,

				Token: accessToken,
			},
		)
	}

	return &Response{

		AccessToken: accessToken,

		RefreshToken: refreshToken,

		AccessExpiresAt: accessClaims.ExpiresAt,

		RefreshExpiresAt: tokenClaims.ExpiresAt,

		Claims: claims.Claims{

			UserID: tokenClaims.UserID,

			SessionID: tokenClaims.SessionID,

			Username: tokenClaims.Username,

			Email: tokenClaims.Email,

			Role: tokenClaims.Role,

			Permissions: tokenClaims.Permissions,
		},
	}, nil
}
