package ctx

import (
	"errors"

	"github.com/qwerius/authgoblue/session"

	"github.com/gofiber/fiber/v3"
)

var errSessionNotFound = errors.New(
	"github.com/qwerius/authgoblue: session not found",
)

const sessionKey = "github.com/qwerius/authgoblue_session"

func (s *Service) SetSession(
	c fiber.Ctx,
	value session.Session,
) {

	c.Locals(
		sessionKey,
		value,
	)
}

func (s *Service) Session(
	c fiber.Ctx,
) (session.Session, error) {

	value := c.Locals(
		sessionKey,
	)

	result, ok := value.(session.Session)

	if !ok {

		return session.Session{}, errSessionNotFound
	}

	return result, nil
}

// SessionID mengambil ID session aktif
func (s *Service) SessionID(
	c fiber.Ctx,
) (string, error) {

	sess, err :=
		s.Session(c)

	if err != nil {

		return "", err
	}

	return sess.ID, nil
}
