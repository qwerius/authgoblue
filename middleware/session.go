package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) RequireSession() fiber.Handler {

	return func(c fiber.Ctx) error {

		authClaims, err := s.context.Claims(c)

		if err != nil {
			return fiber.ErrUnauthorized
		}

		if authClaims.SessionID == "" {
			return fiber.ErrUnauthorized
		}

		if s.session == nil {
			return fiber.ErrUnauthorized
		}

		sess, err := s.session.Get(
			authClaims.SessionID,
		)

		if err != nil {
			return fiber.ErrUnauthorized
		}

		if sess.Revoked {
			return fiber.ErrUnauthorized
		}

		if time.Now().After(
			sess.ExpiresAt,
		) {
			return fiber.ErrUnauthorized
		}

		s.context.SetSession(
			c,
			sess,
		)

		return c.Next()
	}
}
