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
	email string,
	password string,
) (*auth.User, error) {

	return &auth.User{
		ID:       "00000000-0000-0000-0000-000000000001",
		Username: "admin",
		Email:    email,
		Role:     "guest",
	}, nil
}

func (m *mockProvider) Register(
	ctx context.Context,
	input any,
) (*auth.User, error) {

	return &auth.User{
		ID: "00000000-0000-0000-0000-000000000001",
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
		ID: "00000000-0000-0000-0000-000000000001",
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

	agb, err := authgoblue.New(
		authgoblue.Config{
			Secret:   "secret",
			Issuer:   "test",
			Provider: &mockProvider{},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if agb.Client == nil {
		t.Fatal("auth client should not be nil")
	}

	result, err := agb.Client.
		Login().
		Execute(
			context.Background(),
			login.Request{
				Email:    "a@gmail.com",
				Password: "password",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.Result.Email != "a@gmail.com" {
		t.Fatal("invalid email")
	}

	if result.Result.Role != "guest" {
		t.Fatal("invalid role")
	}

	if result.Result.UserID.String() == "" {
		t.Fatal("invalid user id")
	}
}
