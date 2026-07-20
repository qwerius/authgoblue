package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var (
	errTokenNotFound      = errors.New("github.com/qwerius/authgoblue: token not found")
	errMissingHeader      = errors.New("missing header")
	errInvalidAuthFormat  = errors.New("invalid authorization format")
	errInvalidTokenPrefix = errors.New("invalid token prefix")
	errCookieTokenMissing = errors.New("cookie token missing")
)

func (s *Service) extractToken(
	c fiber.Ctx,
) (string, error) {

	if token, err := s.extractHeaderToken(c); err == nil {
		return token, nil
	}

	if s.cookie {

		if token, err := s.extractCookieToken(c); err == nil {
			return token, nil
		}
	}

	return "", errTokenNotFound
}

func (s *Service) extractHeaderToken(
	c fiber.Ctx,
) (string, error) {

	header := c.Get(
		s.header,
	)

	if header == "" {
		return "", errMissingHeader
	}

	prefixLen := len(s.prefix)
	if len(header) <= prefixLen+1 {
		return "", errInvalidAuthFormat
	}

	if header[prefixLen] != ' ' {
		return "", errInvalidAuthFormat
	}

	// Fast case-insensitive prefix check for ASCII
	for i := 0; i < prefixLen; i++ {
		c1 := header[i]
		c2 := s.prefix[i]

		if c1 != c2 {
			// Try case-insensitive match only if needed
			if c1 >= 'A' && c1 <= 'Z' {
				c1 = c1 + 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 = c2 + 32
			}
			if c1 != c2 {
				return "", errInvalidTokenPrefix
			}
		}
	}

	return header[prefixLen+1:], nil
}

func (s *Service) extractCookieToken(
	c fiber.Ctx,
) (string, error) {

	value := c.Cookies(
		s.cookieName,
	)

	if value == "" {
		return "", errCookieTokenMissing
	}

	return value, nil
}
