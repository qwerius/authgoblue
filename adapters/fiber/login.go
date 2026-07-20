package fiberadapter

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/qwerius/authgoblue/login"
)

type LoginHandler struct {
	service *login.Service
}

func NewLoginHandler(
	service *login.Service,
) *LoginHandler {

	return &LoginHandler{
		service: service,
	}
}

func (h *LoginHandler) Handle(
	c fiber.Ctx,
) error {

	var req login.Request

	if err := c.Bind().Body(&req); err != nil {

		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": "invalid request body",
			},
		)
	}

	result, err :=
		h.service.Login(
			c.Context(),
			req,
		)

	if err != nil {

		switch {

		case errors.Is(
			err,
			login.ErrInvalidCredentials,
		):

			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)

		default:

			return c.Status(
				fiber.StatusInternalServerError,
			).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
	}

	return c.JSON(result)
}
