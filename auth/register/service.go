package register

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

	user, err := s.provider.Register(
		ctx,
		req,
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		User: user,
	}, nil
}
