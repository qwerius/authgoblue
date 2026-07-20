package header

import "strings"

func ExtractBearer(
	value string,
	prefix string,
) (string, bool) {

	if value == "" {
		return "", false
	}

	prefixLen := len(prefix)
	if len(value) <= prefixLen+1 {
		return "", false
	}

	if !strings.EqualFold(
		value[:prefixLen],
		prefix,
	) {
		return "", false
	}

	if value[prefixLen] != ' ' {
		return "", false
	}

	token := strings.TrimSpace(
		value[prefixLen+1:],
	)

	if token == "" {
		return "", false
	}

	return token, true
}
