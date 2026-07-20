package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func (s *Service) RequireAuth() fiber.Handler {

	return s.RequireAuth_With(
		Options{},
	)
}

func (s *Service) RequireAuth_With(
	options Options,
) fiber.Handler {

	return func(c fiber.Ctx) error {

		tokenString, err :=
			s.extractToken(c)

		if err != nil {
			return fiber.ErrUnauthorized
		}

		authClaims, err :=
			s.token.ParseAccessToken(
				tokenString,
			)

		if err != nil {
			return fiber.ErrUnauthorized
		}

		if !options.SkipValidation {

			err =
				s.token.ValidateAccessToken(
					authClaims,
				)

			if err != nil {
				return fiber.ErrUnauthorized
			}
		}

		if !options.SkipSessionCheck {

			err =
				s.session.CheckSession(
					authClaims.SessionID,
				)

			if err != nil {
				return fiber.ErrUnauthorized
			}
		}

		// semua lolos
		s.context.SetClaims(
			c,
			authClaims,
		)

		return c.Next()
	}
}
