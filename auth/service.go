package auth

type Service struct {
	provider Provider
}

func New(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) Provider() Provider {
	return s.provider
}
