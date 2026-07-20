package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func (s *Service) RequireRole(
	roleName string,
) fiber.Handler {

	return func(c fiber.Ctx) error {

		if err := s.role.Check(
			c,
			roleName,
		); err != nil {

			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
