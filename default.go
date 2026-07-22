package authgoblue

import (
	"time"

	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
)

const (
	defaultIssuer = "authgoblue"
)

func applyDefaults(config Config) Config {

	if config.Issuer == "" {
		config.Issuer = defaultIssuer
	}

	if config.AccessTokenTTL == 0 {
		config.AccessTokenTTL = 15 * time.Minute
	}

	if config.RefreshTokenTTL == 0 {
		config.RefreshTokenTTL = 7 * 24 * time.Hour
	}

	if config.Header == "" {
		config.Header = "Authorization"
	}

	if config.Prefix == "" {
		config.Prefix = "Bearer"
	}

	if config.AccessCookieName == "" {
		config.AccessCookieName = "access_token"
	}

	if config.RefreshCookieName == "" {
		config.RefreshCookieName = "refresh_token"
	}

	if config.SessionStore == nil {
		config.SessionStore = session.NewMemoryStore()
	}

	if config.RevokeStore == nil {
		config.RevokeStore = revoke.NewMemoryStore()
	}

	if config.MaxSessions == 0 {
		config.MaxSessions = 5
	}

	return config
}
