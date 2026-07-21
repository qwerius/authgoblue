package login

import (
	"context"

	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/claims"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Service struct {
	provider auth.Provider

	token *token.Service

	session *session.Service

	hooks *hooks.Registry

	maxSessions int
}

func New(
	provider auth.Provider,
	tokenService *token.Service,
	sessionService *session.Service,
	hookRegistry *hooks.Registry,
	maxSessions int,
) *Service {

	return &Service{

		provider: provider,

		token: tokenService,

		session: sessionService,

		hooks: hookRegistry,

		maxSessions: maxSessions,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	user, err := s.provider.Authenticate(
		ctx,
		req.Email,
		req.Password,
	)

	if err != nil {
		return nil, err
	}

	if user.Role == "" {
		user.Role = "guest"
	}

	var sess session.Session

	if req.DeviceID != "" ||
		req.DeviceName != "" ||
		req.Platform != "" ||
		req.IPAddress != "" ||
		req.UserAgent != "" {

		sess, err = s.session.CreateWithDevice(
			user.ID,
			req.DeviceID,
			req.DeviceName,
			req.Platform,
			req.IPAddress,
			req.UserAgent,
		)

	} else {

		sess, err = s.session.Create(
			user.ID,
		)
	}

	if err != nil {
		return nil, err
	}

	// Limit active sessions
	err = s.session.EnforceLimit(
		user.ID,
		s.maxSessions,
	)

	if err != nil {
		return nil, err
	}

	authClaims := claims.Claims{

		UserID: user.ID,

		Username: user.Username,

		Email: user.Email,

		Role: user.Role,

		Permissions: user.Permissions,

		SessionID: sess.ID,
	}

	accessToken, err := s.token.GenerateAccessToken(
		authClaims,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.token.GenerateRefreshToken(
		authClaims,
	)

	if err != nil {
		return nil, err
	}

	result := &auth.AuthResult{

		User: user,

		Claims: authClaims,

		AccessToken: accessToken,

		RefreshToken: refreshToken,
	}

	if s.hooks != nil {

		err = s.hooks.Fire(
			ctx,
			hooks.EventAfterLogin,
			hooks.Payload{

				UserID: user.ID,

				SessionID: sess.ID,

				Token: accessToken,

				Metadata: map[string]any{
					"username": user.Username,
					"email":    user.Email,
					"platform": req.Platform,
					"device":   req.DeviceName,
				},
			},
		)

		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Result: result,
	}, nil
}
