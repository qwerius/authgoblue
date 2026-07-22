package client

import (
	"github.com/qwerius/authgoblue/auth"

	"github.com/qwerius/authgoblue/auth/forgotpassword"
	"github.com/qwerius/authgoblue/auth/login"
	"github.com/qwerius/authgoblue/auth/logout"
	"github.com/qwerius/authgoblue/auth/refresh"
	"github.com/qwerius/authgoblue/auth/register"
	"github.com/qwerius/authgoblue/auth/resetpassword"
	"github.com/qwerius/authgoblue/auth/verifyemail"
	"github.com/qwerius/authgoblue/ctx"
	"github.com/qwerius/authgoblue/hooks"
	"github.com/qwerius/authgoblue/password"
	coreRefresh "github.com/qwerius/authgoblue/refresh"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/token"
)

type Client struct {
	provider auth.Provider

	token   *token.Service
	session *session.Service
	revoke  *revoke.Service
	ctx     *ctx.Service

	hooks *hooks.Registry

	maxSessions int
	password    *password.Service
}

func New(
	provider auth.Provider,
	tokenService *token.Service,
	sessionService *session.Service,
	revokeService *revoke.Service,
	ctxService *ctx.Service,
	hookRegistry *hooks.Registry,
	maxSessions int,
	password *password.Service,
) *Client {

	return &Client{

		provider: provider,

		token: tokenService,

		session: sessionService,

		revoke: revokeService,

		ctx: ctxService,

		hooks: hookRegistry,

		maxSessions: maxSessions,
		password:    password,
	}
}

func (c *Client) Login() *login.Service {

	return login.New(
		c.provider,
		c.token,
		c.session,
		c.hooks,
		c.maxSessions,
	)
}

func (c *Client) Logout() *logout.Service {

	return logout.New(
		c.provider,
		c.token,
		c.session,
		c.revoke,
		c.hooks,
		c.ctx,
	)
}

func (c *Client) Refresh() *refresh.Service {

	rotateService :=
		coreRefresh.NewService(
			c.token,
			c.revoke,
			c.session,
		)

	return refresh.New(
		rotateService,
		c.hooks,
	)
}

func (c *Client) Register() *register.Service {

	return register.New(
		c.provider,
	)
}

func (c *Client) ForgotPassword() *forgotpassword.Service {

	return forgotpassword.New(
		c.provider,
	)
}

func (c *Client) ResetPassword() *resetpassword.Service {

	return resetpassword.New(
		c.provider,
		c.password,
	)
}

func (c *Client) VerifyEmail() *verifyemail.Service {

	return verifyemail.New(
		c.provider,
	)
}

func (c *Client) Session() *session.Service {

	return c.session
}

func (c *Client) Hooks() *hooks.Registry {

	return c.hooks
}

func (c *Client) MaxSessions() int {

	return c.maxSessions
}
