package authgoblue

import (
	"time"

	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
)

type Config struct {

	// Secret key untuk signing JWT
	Secret string

	// Nama aplikasi/service
	Issuer string

	// Durasi access token
	AccessTokenTTL time.Duration

	// Durasi refresh token
	RefreshTokenTTL time.Duration

	Header string

	Prefix string

	Cookie bool

	CookieName string

	// Auth provider
	Provider auth.Provider

	// Session storage
	// default MemoryStore
	SessionStore session.Store

	// Maximum active sessions per user
	// default 5
	MaxSessions int

	// Token revoke storage
	// default MemoryStore
	RevokeStore revoke.Store
}

func DefaultConfig() Config {

	return Config{

		AccessTokenTTL: 15 * time.Minute,

		RefreshTokenTTL: 7 * 24 * time.Hour,

		Header: "Authorization",

		Prefix: "Bearer",

		Cookie: false,

		CookieName: "github.com/qwerius/authgoblue_token",

		MaxSessions: 5,
	}
}
