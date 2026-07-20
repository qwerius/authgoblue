package role

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var errRoleForbidden = errors.New("authgoblue: role forbidden")

func (s *Service) Get(
	c fiber.Ctx,
) (string, error) {

	return s.context.Role(c)
}

func (s *Service) Check(
	c fiber.Ctx,
	requiredRole string,
) error {

	currentRole, err := s.Get(c)

	if err != nil {
		return err
	}

	if currentRole != requiredRole {
		return errRoleForbidden
	}

	return nil
}

func (s *Service) Has(
	c fiber.Ctx,
	requiredRole string,
) (bool, error) {

	currentRole, err := s.Get(c)

	if err != nil {
		return false, err
	}

	return currentRole == requiredRole, nil
}
