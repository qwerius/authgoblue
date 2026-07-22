package resetpassword

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
	passwordService *password.Service,
) *Service {
	return &Service{
		provider: provider,
		password: passwordService,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	req Request,
) (*Response, error) {

	user, err := s.provider.ValidateResetToken(
		ctx,
		req.Token,
	)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.password.Hash(
		req.NewPassword,
	)
	if err != nil {
		return nil, err
	}

	err = s.provider.UpdatePassword(
		ctx,
		user.ID,
		hashedPassword,
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		Success: true,
	}, nil
}
