package util

import "strings"

func Empty(
	value string,
) bool {

	return strings.TrimSpace(
		value,
	) == ""
}

func Normalize(
	value string,
) string {

	return strings.TrimSpace(
		strings.ToLower(
			value,
		),
	)
}
