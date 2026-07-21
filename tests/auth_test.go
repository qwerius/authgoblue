package authgoblue_test

import (
	"context"
	"testing"

	"github.com/qwerius/authgoblue"
	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/auth/login"
)

type mockProvider struct{}

func (m *mockProvider) Authenticate(
	ctx context.Context,
	username string,
	password string,
) (*auth.User, error) {

	return &auth.User{
		ID:       "1",
		Username: username,
		Email:    "test@example.com",
		Role:     "guest",
	}, nil
}

func (m *mockProvider) Register(
	ctx context.Context,
	input any,
) (*auth.User, error) {

	return &auth.User{
		ID: "1",
	}, nil
}

func (m *mockProvider) FindByID(
	ctx context.Context,
	id string,
) (*auth.User, error) {

	return &auth.User{
		ID: id,
	}, nil
}

func (m *mockProvider) FindByEmail(
	ctx context.Context,
	email string,
) (*auth.User, error) {

	return &auth.User{
		Email: email,
	}, nil
}

func (m *mockProvider) SaveResetToken(
	ctx context.Context,
	userID string,
	token string,
) error {

	return nil
}

func (m *mockProvider) ValidateResetToken(
	ctx context.Context,
	token string,
) (*auth.User, error) {

	return &auth.User{
		ID: "1",
	}, nil
}

func (m *mockProvider) UpdatePassword(
	ctx context.Context,
	userID string,
	hashedPassword string,
) error {

	return nil
}

func (m *mockProvider) VerifyEmail(
	ctx context.Context,
	token string,
) error {

	return nil
}

func TestAuthClientLogin(t *testing.T) {

	provider := &mockProvider{}

	agb := authgoblue.New(
		authgoblue.Config{
			Secret:   "secret",
			Issuer:   "test",
			Provider: provider,
		},
	)

	if agb.Client == nil {
		t.Fatal("auth client should not be nil")
	}

	result, err := agb.Client.
		Login().
		Execute(
			context.Background(),
			login.Request{
				Username: "admin",
				Password: "password",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.Result.User.Username != "admin" {
		t.Fatal("invalid user")
	}
}
