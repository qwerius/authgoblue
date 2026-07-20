package providers

import "context"

// User represents identity data needed by authentication.
type User struct {
	ID string

	Username string

	Email string

	PasswordHash string

	Role string

	Permissions []string
}

// Provider is implemented by application.
//
// Example:
// PostgreSQL user repository,
// LDAP,
// external API,
// etc.
type Provider interface {
	FindByIdentifier(
		ctx context.Context,
		identifier string,
	) (*User, error)
}
