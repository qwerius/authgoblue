package authgoblue

import (
	"time"

	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/revoke"
	"github.com/qwerius/authgoblue/session"
)

const (
	DefaultIssuer = "authgoblue"
)

type Config struct {

	// Required:
	// Secret key untuk signing JWT.
	// Consumer harus mengganti dengan secret sendiri.
	Secret string

	// Default: authgoblue
	// Nama aplikasi/service penerbit token.
	Issuer string

	// Default: 15 menit
	// Durasi access token.
	AccessTokenTTL time.Duration

	// Default: 7 hari
	// Durasi refresh token.
	RefreshTokenTTL time.Duration

	// Default: Authorization
	Header string

	// Default: Bearer
	Prefix string

	// Default: false
	// Jika true, token disimpan melalui cookie.
	Cookie bool

	// Default: access_token
	AccessCookieName string

	// Default: refresh_token
	RefreshCookieName string

	// Required:
	// Provider autentikasi user.
	Provider auth.Provider

	// Optional:
	// Default: MemoryStore.
	SessionStore session.Store

	// Default: 5
	// Maximum session aktif per user.
	MaxSessions int

	// Optional:
	// Default: MemoryStore.
	RevokeStore revoke.Store
}
