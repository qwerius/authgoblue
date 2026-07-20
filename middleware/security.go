package middleware

import (
	"time"

	"github.com/qwerius/authgoblue/claims"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) checkSession(
	c claims.Claims,
) error {

	// Session tidak digunakan
	if s.session == nil {
		return nil
	}

	// Token tanpa SessionID
	if c.SessionID == "" {
		return fiber.ErrUnauthorized
	}

	sess, err := s.session.Get(
		c.SessionID,
	)

	if err != nil {
		// JWT valid tapi session tidak ditemukan
		return fiber.ErrUnauthorized
	}

	if sess.Revoked {
		return fiber.ErrUnauthorized
	}

	if !sess.ExpiresAt.IsZero() &&
		time.Now().After(
			sess.ExpiresAt,
		) {
		return fiber.ErrUnauthorized
	}

	return nil
}
