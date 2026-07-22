package authgoblue_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/auth/forgotpassword"
	"github.com/qwerius/authgoblue/password"
)

type forgotPasswordMockProvider struct {
	findEmailCalled bool
	saveTokenCalled bool

	userID string
	token  string
}

func (m *forgotPasswordMockProvider) FindByEmail(
	ctx context.Context,
	email string,
) (*auth.User, error) {

	m.findEmailCalled = true

	return &auth.User{
		ID:    m.userID,
		Email: email,
	}, nil
}

func (m *forgotPasswordMockProvider) SaveResetToken(
	ctx context.Context,
	userID string,
	token string,
) error {

	m.saveTokenCalled = true
	m.token = token

	return nil
}

func (m *forgotPasswordMockProvider) ValidateResetToken(
	ctx context.Context,
	token string,
) (*auth.User, error) {

	return nil, nil
}

func (m *forgotPasswordMockProvider) UpdatePassword(
	ctx context.Context,
	userID string,
	hashedPassword string,
) error {

	return nil
}

func (m *forgotPasswordMockProvider) Authenticate(
	ctx context.Context,
	email string,
	password string,
) (*auth.User, error) {

	return nil, nil
}

func (m *forgotPasswordMockProvider) Register(
	ctx context.Context,
	input any,
) (*auth.User, error) {

	return nil, nil
}

func (m *forgotPasswordMockProvider) FindByID(
	ctx context.Context,
	id string,
) (*auth.User, error) {

	return nil, nil
}

func (m *forgotPasswordMockProvider) VerifyEmail(
	ctx context.Context,
	token string,
) error {

	return nil
}

func TestForgotPassword_Success(
	t *testing.T,
) {

	provider := &forgotPasswordMockProvider{
		userID: "user-123",
	}

	passwordService := password.NewService()

	service := forgotpassword.New(
		provider,
		passwordService,
	)

	response, err := service.Execute(
		context.Background(),
		forgotpassword.Request{
			Email: "user@test.com",
		},
	)

	require.NoError(
		t,
		err,
	)

	assert.True(
		t,
		response.Success,
	)

	assert.NotEmpty(
		t,
		response.Token,
	)

	assert.True(
		t,
		provider.findEmailCalled,
	)

	assert.True(
		t,
		provider.saveTokenCalled,
	)

	assert.Equal(
		t,
		response.Token,
		provider.token,
	)
}
