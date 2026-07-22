package refresh

import (
	"context"
	"time"

	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	coreRefresh "github.com/qwerius/authgoblue/refresh"
)

type Service struct {
	refresh *coreRefresh.Service

	hooks *hooks.Registry
}

func New(
	refreshService *coreRefresh.Service,
	hookRegistry *hooks.Registry,
) *Service {

	return &Service{
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
		accessExpiresAt,
		refreshExpiresAt,
		err :=
		s.refresh.Rotate(
			req.RefreshToken,
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

	now := time.Now().Unix()

	return &Response{

		AccessToken: accessToken,

		RefreshToken: refreshToken,

		AccessExpiresAt: accessExpiresAt,
		AccessExpiresIn: accessExpiresAt - now,

		RefreshExpiresAt: refreshExpiresAt,
		RefreshExpiresIn: refreshExpiresAt - now,

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
