package token

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/qwerius/authgoblue/claims"
)

const jwtHeader = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

var errMarshalClaims = errors.New("github.com/qwerius/authgoblue: failed to marshal claims")

func (s *Service) encodeJWT(
	c claims.Claims,
) (string, error) {

	payload, err := json.Marshal(c)

	if err != nil {
		return "", errMarshalClaims
	}

	payloadSegment := encodeSegment(
		payload,
	)

	unsignedToken :=
		jwtHeader + "." + payloadSegment

	signature := signHS256(
		[]byte(unsignedToken),
		s.secret,
	)

	signatureSegment := encodeSegment(
		signature,
	)

	var builder strings.Builder

	builder.Grow(
		len(jwtHeader) +
			1 +
			len(payloadSegment) +
			1 +
			len(signatureSegment),
	)

	builder.WriteString(jwtHeader)
	builder.WriteByte('.')
	builder.WriteString(payloadSegment)
	builder.WriteByte('.')
	builder.WriteString(signatureSegment)

	return builder.String(), nil
}
