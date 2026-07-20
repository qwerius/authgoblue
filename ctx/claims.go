package ctx

import (
	"authgoblue/claims"
	"errors"

	"github.com/gofiber/fiber/v3"
)

var errClaimsNotFound = errors.New("authgoblue: claims not found")

const claimsKey = "authgoblue_claims"

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
