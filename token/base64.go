package token

import (
	"encoding/base64"
)

var jwtEncoding = base64.RawURLEncoding

func encodeSegment(
	data []byte,
) string {

	return jwtEncoding.EncodeToString(
		data,
	)
}

func decodeSegment(
	segment string,
) ([]byte, error) {

	return jwtEncoding.DecodeString(
		segment,
	)
}
