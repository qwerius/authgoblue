package permission

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var errPermissionForbidden = errors.New("github.com/qwerius/authgoblue: permission forbidden")

func (s *Service) Get(
	c fiber.Ctx,
) ([]string, error) {

	return s.context.Permissions(c)
}

func (s *Service) Check(
	c fiber.Ctx,
	requiredPermission string,
) error {

	allowed, err := s.Has(
		c,
		requiredPermission,
	)

	if err != nil {
		return err
	}

	if !allowed {
		return errPermissionForbidden
	}

	return nil
}

func (s *Service) Has(
	c fiber.Ctx,
	requiredPermission string,
) (bool, error) {

	permissions, err := s.Get(c)

	if err != nil {
		return false, err
	}

	for _, permission := range permissions {

		if permission == requiredPermission {
			return true, nil
		}
	}

	return false, nil
}
