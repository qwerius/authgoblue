package password

import (
	"crypto/rand"
	"encoding/hex"

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

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	if err != nil {
		return ErrInvalidPassword
	}

	return nil
}

func (s *Service) GenerateResetToken() string {

	bytes := make(
		[]byte,
		32,
	)

	_, err := rand.Read(
		bytes,
	)

	if err != nil {
		return ""
	}

	return hex.EncodeToString(
		bytes,
	)
}
