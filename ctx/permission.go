package ctx

import "github.com/gofiber/fiber/v3"

func (s *Service) Permissions(
	c fiber.Ctx,
) ([]string, error) {

	authClaims, err := s.Claims(c)

	if err != nil {
		return nil, err
	}

	return authClaims.Permissions, nil
}

func (s *Service) HasPermission(
	c fiber.Ctx,
	permission string,
) (bool, error) {

	permissions, err := s.Permissions(c)

	if err != nil {
		return false, err
	}

	for _, item := range permissions {

		if item == permission {
			return true, nil
		}
	}

	return false, nil
}
