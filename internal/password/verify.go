package password

import "golang.org/x/crypto/bcrypt"

func (s *Service) Verify(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
