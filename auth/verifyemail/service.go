package verifyemail

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

	err := s.provider.VerifyEmail(
		ctx,
		req.Token,
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		Success: true,
	}, nil
}
