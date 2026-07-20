package authgoblue

import "github.com/gofiber/fiber/v3"

func (a *AuthGoBlue) UserID(
	c fiber.Ctx,
) (string, error) {

	return a.Context.UserID(c)
}

func (a *AuthGoBlue) Username(
	c fiber.Ctx,
) (string, error) {

	return a.Context.Username(c)
}

func (a *AuthGoBlue) Email(
	c fiber.Ctx,
) (string, error) {

	return a.Context.Email(c)
}

func (a *AuthGoBlue) Role(
	c fiber.Ctx,
) (string, error) {

	return a.Context.Role(c)
}

func (a *AuthGoBlue) Permissions(
	c fiber.Ctx,
) ([]string, error) {

	return a.Context.Permissions(c)
}

func (a *AuthGoBlue) HasPermission(
	c fiber.Ctx,
	permission string,
) (bool, error) {

	return a.Context.HasPermission(
		c,
		permission,
	)
}
