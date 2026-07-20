package ctx

import "github.com/gofiber/fiber/v3"

func (s *Service) UserID(
	c fiber.Ctx,
) (string, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return "", err
	}

	return authClaims.UserID, nil
}

func (s *Service) Username(
	c fiber.Ctx,
) (string, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return "", err
	}

	return authClaims.Username, nil
}

func (s *Service) Email(
	c fiber.Ctx,
) (string, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return "", err
	}

	return authClaims.Email, nil
}
