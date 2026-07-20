package ctx

import "github.com/gofiber/fiber/v3"

func (s *Service) Role(
	c fiber.Ctx,
) (string, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return "", err
	}

	return authClaims.Role, nil
}
