package logout

import (
	"context"
	"time"

	"github.com/qwerius/authgoblue/auth"
	authctx "github.com/qwerius/authgoblue/ctx"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	provider auth.Provider

	token *token.Service

	session *session.Service

	revoke *revoke.Service

	hooks *hooks.Registry

	ctx *authctx.Service
}

func New(
	provider auth.Provider,
	tokenService *token.Service,
	sessionService *session.Service,
	revokeService *revoke.Service,
	hookRegistry *hooks.Registry,
	ctxService *authctx.Service,
) *Service {

	return &Service{
		provider: provider,

		token: tokenService,

		session: sessionService,

		revoke: revokeService,

		hooks: hookRegistry,

		ctx: ctxService,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) error {

	if req.RefreshToken == "" {
		return ErrInvalidRefreshToken
	}

	claims, err :=
		s.token.ParseRefreshToken(
			req.RefreshToken,
		)

	if err != nil {
		return ErrInvalidRefreshToken
	}

	err =
		s.token.ValidateRefreshToken(
			claims,
		)

	if err != nil {
		return ErrInvalidRefreshToken
	}

	// revoke refresh token lama
	if s.revoke != nil {

		err =
			s.revoke.Revoke(
				claims.TokenID,
				time.Unix(
					claims.ExpiresAt,
					0,
				),
			)

		if err != nil {
			return ErrRevokeToken
		}
	}

	// hapus session
	if s.session != nil &&
		claims.SessionID != "" {

		err =
			s.session.Revoke(
				claims.SessionID,
			)

		if err != nil {
			return ErrDeleteSession
		}
	}

	// hook logout
	if s.hooks != nil {

		err =
			s.hooks.Fire(
				ctx,
				hooks.EventAfterLogout,
				hooks.Payload{

					UserID: claims.UserID,

					SessionID: claims.SessionID,

					Metadata: map[string]any{
						"email": claims.Email,
					},
				},
			)

		if err != nil {
			return ErrHookExecution
		}
	}

	return nil
}
