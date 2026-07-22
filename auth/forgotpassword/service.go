package forgotpassword

import (
	"context"

	"github.com/qwerius/authgoblue/auth"
	"github.com/qwerius/authgoblue/password"
)

type Service struct {
	provider auth.Provider
	password *password.Service
}

func New(
	provider auth.Provider,
	passwordService *password.Service) *Service {
	return &Service{
		provider: provider,
		password: passwordService,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	user, err := s.provider.FindByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return nil, err
	}

	resetToken := s.password.GenerateResetToken()

	err = s.provider.SaveResetToken(
		ctx,
		user.ID,
		resetToken,
	)

	if err != nil {
		return nil, err
	}

	return &Response{
		Success: true,
		Token:   resetToken, // sementara untuk testing
	}, nil
}
