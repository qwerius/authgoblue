package password

import (
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	cost int
}

func NewService() *Service {

	return &Service{
		cost: bcrypt.DefaultCost,
	}
}

func NewServiceWithCost(
	cost int,
) *Service {

	return &Service{
		cost: cost,
	}
}

func (s *Service) Hash(
	password string,
) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		s.cost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *Service) Compare(
	password string,
	hash string,
) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
