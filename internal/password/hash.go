package password

import "golang.org/x/crypto/bcrypt"

func (s *Service) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		s.cost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
