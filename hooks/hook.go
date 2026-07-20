package hooks

import "context"

// Event represents an auth lifecycle event.
type Event string

const (

	// Login events
	//EventBeforeLogin Event = "before_login"
	EventAfterLogin Event = "after_login"

	// Token events
	//EventBeforeTokenGenerate Event = "before_token_generate"
	//EventAfterTokenGenerate  Event = "after_token_generate"

	// Refresh events
	//EventBeforeRefresh Event = "before_refresh"
	//EventAfterRefresh  Event = "after_refresh"

	// Session events
	//EventSessionCreated Event = "session_created"
	//EventSessionRevoked Event = "session_revoked"

	// Logout events
	//EventBeforeLogout Event = "before_logout"
	//EventAfterLogout  Event = "after_logout"
)

// Payload contains data passed to hook handlers.
type Payload struct {
	UserID string

	SessionID string

	Token string

	Metadata map[string]any
}

// Handler executes when an event occurs.
type Handler func(
	ctx context.Context,
	payload Payload,
) error
