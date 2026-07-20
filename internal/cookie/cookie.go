package cookie

import (
	"github.com/gofiber/fiber/v3"
)

func Get(
	c fiber.Ctx,
	name string,
) (string, bool) {

	value := c.Cookies(
		name,
	)

	if value == "" {

		return "", false
	}

	return value, true
}
