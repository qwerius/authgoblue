package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func (s *Service) RequirePermission(
	permissionName string,
) fiber.Handler {

	return func(c fiber.Ctx) error {

		if err := s.permission.Check(
			c,
			permissionName,
		); err != nil {

			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
