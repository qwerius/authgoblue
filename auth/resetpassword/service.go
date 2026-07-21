package resetpassword

import (
	"context"

	"github.com/qwerius/authgoblue/auth"
)

type Service struct {
	provider auth.Provider
}

func New(provider auth.Provider) *Service {
	return &Service{
		provider: provider,
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

	// Password hashing nanti masuk layer password
	hashedPassword := req.NewPassword

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
