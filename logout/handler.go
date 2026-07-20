package logout

import (
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service

	cookieName string
}

func NewHandler(
	service *Service,
	cookieName string,
) *Handler {

	return &Handler{

		service: service,

		cookieName: cookieName,
	}
}

func (h *Handler) Logout() fiber.Handler {

	return func(c fiber.Ctx) error {

		token :=
			extractToken(
				c,
				h.cookieName,
			)

		err :=
			h.service.Logout(
				token,
			)

		if err != nil {

			return fiber.ErrUnauthorized
		}

		c.Cookie(
			&fiber.Cookie{

				Name: h.cookieName,

				Value: "",

				MaxAge: -1,

				HTTPOnly: true,
			},
		)

		return c.SendStatus(
			fiber.StatusNoContent,
		)
	}
}

func extractToken(
	c fiber.Ctx,
	cookieName string,
) string {

	auth :=
		c.Get(
			"Authorization",
		)

	if auth != "" {

		const prefix = "Bearer "

		if len(auth) > len(prefix) &&
			auth[:len(prefix)] == prefix {

			return auth[len(prefix):]
		}
	}

	return c.Cookies(
		cookieName,
	)
}
