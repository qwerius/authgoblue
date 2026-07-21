package authgoblue

import (
	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/auth/client"

	"github.com/qwerius/authgoblue/ctx"
	"github.com/qwerius/authgoblue/hooks"

	"github.com/qwerius/authgoblue/middleware"
	"github.com/qwerius/authgoblue/password"
	"github.com/qwerius/authgoblue/permission"

	"github.com/qwerius/authgoblue/refresh"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/role"
	"github.com/qwerius/authgoblue/session"
	"github.com/qwerius/authgoblue/storage"
	"github.com/qwerius/authgoblue/token"
)

type AuthGoBlue struct {
	config Config

	Token       *token.Service
	Password    *password.Service
	Context     *ctx.Service
	RoleService *role.Service
	Permission  *permission.Service
	Revoke      *revoke.Service
	Session     *session.Service
	Middleware  *middleware.Service

	Storage *storage.Registry
	Hooks   *hooks.Registry

	Refresh *refresh.Service

	Auth   *auth.Service
	Client *client.Client
}

func New(config Config) *AuthGoBlue {

	config = applyDefaults(config)

	agb := &AuthGoBlue{
		config: config,
	}

	// Context
	agb.Context = ctx.NewService()

	// Token
	agb.Token = token.NewService(
		config.Secret,
		config.Issuer,
		config.AccessTokenTTL,
		config.RefreshTokenTTL,
	)

	// Password
	agb.Password = password.NewService()

	// Role
	agb.RoleService = role.NewService(
		agb.Context,
	)

	// Permission
	agb.Permission = permission.NewService(
		agb.Context,
	)

	// Revoke
	var revokeStore revoke.Store

	if config.RevokeStore != nil {

		revokeStore = config.RevokeStore

	} else {

		revokeStore = revoke.NewMemoryStore()
	}

	agb.Revoke = revoke.NewService(
		revokeStore,
	)

	// Hooks
	agb.Hooks = hooks.NewRegistry()

	// Session
	var sessionStore session.Store

	if config.SessionStore != nil {

		sessionStore = config.SessionStore

	} else {

		sessionStore = session.NewMemoryStore()
	}

	agb.Session = session.NewService(
		sessionStore,
		agb.Hooks,
	)

	// Refresh
	agb.Refresh = refresh.NewService(
		agb.Token,
		agb.Revoke,
		agb.Session,
	)

	// Middleware
	agb.Middleware = middleware.NewService(
		agb.Token,
		agb.Context,
		agb.RoleService,
		agb.Permission,
		agb.Revoke,
		agb.Session,

		config.Header,
		config.Prefix,
		config.Cookie,
		config.CookieName,
	)

	// Storage
	agb.Storage = storage.NewRegistry()

	// Hooks
	agb.Hooks = hooks.NewRegistry()

	// Auth Layer
	if config.Provider != nil {

		agb.Auth = auth.New(
			config.Provider,
		)

		agb.Client = client.New(
			config.Provider,
			agb.Token,
			agb.Session,
			agb.Hooks,
			config.MaxSessions,
		)
	}

	return agb
}

func (a *AuthGoBlue) Config() Config {
	return a.config
}
