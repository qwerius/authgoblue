package forgotpassword

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

	user, err := s.provider.FindByEmail(
		ctx,
		req.Email,
	)
	if err != nil {
		return nil, err
	}

	// Generate reset token nanti diletakkan di layer auth/token
	// sementara hanya simpan token dummy

	err = s.provider.SaveResetToken(
		ctx,
		user.ID,
		"reset-token",
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		Success: true,
	}, nil
}
