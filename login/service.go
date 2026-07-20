package login

import (
	"context"

	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/password"
	"github.com/qwerius/authgoblue/providers"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	provider providers.Provider

	password *password.Service

	token *token.Service

	session *session.Service

	hooks *hooks.Registry
}

func NewService(
	provider providers.Provider,
	passwordService *password.Service,
	tokenService *token.Service,
	sessionService *session.Service,
	hooksRegistry *hooks.Registry,
) *Service {

	return &Service{

		provider: provider,

		password: passwordService,

		token: tokenService,

		session: sessionService,

		hooks: hooksRegistry,
	}
}

type Request struct {
	Identifier string

	Password string
}

type Result struct {
	User *providers.User

	Session session.Session

	AccessToken string

	RefreshToken string
}

func (s *Service) Login(
	ctx context.Context,
	req Request,
) (*Result, error) {

	user, err :=
		s.provider.FindByIdentifier(
			ctx,
			req.Identifier,
		)

	if err != nil || user == nil {

		return nil, ErrInvalidCredentials
	}

	err =
		s.password.Compare(
			req.Password,
			user.PasswordHash,
		)

	if err != nil {

		return nil, ErrInvalidCredentials
	}

	sess, err :=
		s.session.Create(
			user.ID,
		)

	if err != nil {

		return nil, err
	}

	tokenClaims :=
		claims.Claims{

			UserID: user.ID,

			SessionID: sess.ID,

			Username: user.Username,

			Email: user.Email,

			Role: user.Role,

			Permissions: user.Permissions,
		}

	accessToken, err :=
		s.token.GenerateAccessToken(
			tokenClaims,
		)

	if err != nil {

		return nil, err
	}

	refreshToken, err :=
		s.token.GenerateRefreshToken(
			tokenClaims,
		)

	if err != nil {

		return nil, err
	}

	if s.hooks != nil {

		_ = s.hooks.Fire(
			ctx,
			hooks.EventAfterLogin,
			hooks.Payload{

				UserID: user.ID,

				SessionID: sess.ID,
			},
		)
	}

	return &Result{

		User: user,

		Session: sess,

		AccessToken: accessToken,

		RefreshToken: refreshToken,
	}, nil
}
