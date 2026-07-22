package authgoblue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/auth/resetpassword"
	"github.com/qwerius/authgoblue/password"
)

type resetPasswordMockProvider struct {
	validateCalled bool
	updateCalled   bool

	validateError error
	updateError   error

	user *auth.User
	hash string
}

func (m *resetPasswordMockProvider) ValidateResetToken(
	ctx context.Context,
	token string,
) (*auth.User, error) {

	m.validateCalled = true

	if m.validateError != nil {
		return nil, m.validateError
	}

	return m.user, nil
}

func (m *resetPasswordMockProvider) UpdatePassword(
	ctx context.Context,
	userID string,
	hashedPassword string,
) error {

	m.updateCalled = true
	m.hash = hashedPassword

	if m.updateError != nil {
		return m.updateError
	}

	return nil
}

func (m *resetPasswordMockProvider) SaveResetToken(
	ctx context.Context,
	userID string,
	token string,
) error {
	return nil
}

func (m *resetPasswordMockProvider) Authenticate(
	ctx context.Context,
	email string,
	password string,
) (*auth.User, error) {
	return nil, nil
}

func (m *resetPasswordMockProvider) Register(
	ctx context.Context,
	input any,
) (*auth.User, error) {
	return nil, nil
}

func (m *resetPasswordMockProvider) FindByID(
	ctx context.Context,
	id string,
) (*auth.User, error) {
	return nil, nil
}

func (m *resetPasswordMockProvider) FindByEmail(
	ctx context.Context,
	email string,
) (*auth.User, error) {
	return nil, nil
}

func (m *resetPasswordMockProvider) VerifyEmail(
	ctx context.Context,
	token string,
) error {
	return nil
}

func TestResetPassword_Success(t *testing.T) {

	provider := &resetPasswordMockProvider{
		user: &auth.User{
			ID: "user-123",
		},
	}

	passwordService := password.NewService()

	service := resetpassword.New(
		provider,
		passwordService,
	)

	response, err := service.Execute(
		context.Background(),
		resetpassword.Request{
			Token:       "reset-token",
			NewPassword: "PasswordBaru123",
		},
	)

	require.NoError(t, err)

	assert.True(
		t,
		response.Success,
	)

	assert.True(
		t,
		provider.validateCalled,
	)

	assert.True(
		t,
		provider.updateCalled,
	)

	assert.NotEmpty(
		t,
		provider.hash,
	)

	err = passwordService.Compare(
		"PasswordBaru123",
		provider.hash,
	)

	assert.NoError(
		t,
		err,
	)
}

func TestResetPassword_InvalidToken(t *testing.T) {

	provider := &resetPasswordMockProvider{
		validateError: errors.New("invalid reset token"),
	}

	service := resetpassword.New(
		provider,
		password.NewService(),
	)

	response, err := service.Execute(
		context.Background(),
		resetpassword.Request{
			Token:       "invalid-token",
			NewPassword: "PasswordBaru123",
		},
	)

	require.Error(t, err)

	assert.Nil(
		t,
		response,
	)

	assert.True(
		t,
		provider.validateCalled,
	)

	assert.False(
		t,
		provider.updateCalled,
	)
}

func TestResetPassword_UpdatePasswordFailed(t *testing.T) {

	provider := &resetPasswordMockProvider{
		user: &auth.User{
			ID: "user-123",
		},
		updateError: errors.New("update password failed"),
	}

	service := resetpassword.New(
		provider,
		password.NewService(),
	)

	response, err := service.Execute(
		context.Background(),
		resetpassword.Request{
			Token:       "reset-token",
			NewPassword: "PasswordBaru123",
		},
	)

	require.Error(t, err)

	assert.Nil(
		t,
		response,
	)

	assert.True(
		t,
		provider.validateCalled,
	)

	assert.True(
		t,
		provider.updateCalled,
	)

	assert.NotEmpty(
		t,
		provider.hash,
	)
}
