package ctx

import (
	"errors"

	"github.com/qwerius/authgoblue/claims"

	"github.com/gofiber/fiber/v3"
)

var errClaimsNotFound = errors.New("github.com/qwerius/authgoblue: claims not found")

const claimsKey = "github.com/qwerius/authgoblue_claims"

func (s *Service) SetClaims(
	c fiber.Ctx,
	value claims.Claims,
) {
	c.Locals(
		claimsKey,
		value,
	)
}

func (s *Service) Claims(
	c fiber.Ctx,
) (claims.Claims, error) {

	value := c.Locals(
		claimsKey,
	)

	result, ok := value.(claims.Claims)

	if !ok {
		return claims.Claims{}, errClaimsNotFound
	}

	return result, nil
}
