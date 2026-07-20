package password

import "golang.org/x/crypto/bcrypt"

type Service struct {
	cost int
}

func NewService() *Service {
	return &Service{
		cost: bcrypt.DefaultCost,
	}
}

func NewServiceWithCost(cost int) *Service {
	return &Service{
		cost: cost,
	}
}

func (s *Service) Cost() int {
	return s.cost
}
