package refresh

import (
	"context"

	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	token *token.Service

	session *session.Service

	hooks *hooks.Registry
}

func New(
	tokenService *token.Service,
	sessionService *session.Service,
	hookRegistry *hooks.Registry,
) *Service {

	return &Service{

		token: tokenService,

		session: sessionService,

		hooks: hookRegistry,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	// validate refresh token

	oldClaims, err :=
		s.token.ParseRefreshToken(
			req.RefreshToken,
		)

	if err != nil {
		return nil, err
	}

	// check old session

	err = s.session.CheckSession(
		oldClaims.SessionID,
	)

	if err != nil {
		return nil, err
	}

	// revoke old session

	err = s.session.Revoke(
		oldClaims.SessionID,
	)

	if err != nil {
		return nil, err
	}

	// create new session

	newSession, err :=
		s.session.Create(
			oldClaims.UserID,
		)

	if err != nil {
		return nil, err
	}

	newClaims := claims.Claims{

		UserID: oldClaims.UserID,

		Username: oldClaims.Username,

		Email: oldClaims.Email,

		Role: oldClaims.Role,

		Permissions: oldClaims.Permissions,

		SessionID: newSession.ID,
	}

	accessToken, err :=
		s.token.GenerateAccessToken(
			newClaims,
		)

	if err != nil {
		return nil, err
	}

	refreshToken, err :=
		s.token.GenerateRefreshToken(
			newClaims,
		)

	if err != nil {
		return nil, err
	}

	if s.hooks != nil {

		_ = s.hooks.Fire(
			ctx,
			hooks.EventAfterRefresh,
			hooks.Payload{

				UserID: newClaims.UserID,

				SessionID: newSession.ID,

				Token: accessToken,
			},
		)
	}

	return &Response{

		AccessToken: accessToken,

		RefreshToken: refreshToken,

		Claims: newClaims,
	}, nil
}
