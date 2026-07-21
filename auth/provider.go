package auth

import "context"

type Authenticator interface {
	Authenticate(
		ctx context.Context,
		email string,
		password string,
	) (*User, error)
}

type Registrar interface {
	Register(
		ctx context.Context,
		input any,
	) (*User, error)
}

type UserFinder interface {
	FindByID(
		ctx context.Context,
		id string,
	) (*User, error)

	FindByEmail(
		ctx context.Context,
		email string,
	) (*User, error)
}

type PasswordResetter interface {
	SaveResetToken(
		ctx context.Context,
		userID string,
		token string,
	) error

	ValidateResetToken(
		ctx context.Context,
		token string,
	) (*User, error)

	UpdatePassword(
		ctx context.Context,
		userID string,
		hashedPassword string,
	) error
}

type EmailVerifier interface {
	VerifyEmail(
		ctx context.Context,
		token string,
	) error
}

type Provider interface {
	Authenticator
	Registrar
	UserFinder
	PasswordResetter
	EmailVerifier
}
