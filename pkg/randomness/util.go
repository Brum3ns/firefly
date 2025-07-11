package randomness

import (
	"regexp"
	"strings"
)

var (
	uuidRegex   = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	hexRegex    = regexp.MustCompile(`^[a-f0-9]{32,}$`)
	base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/=]{16,}$`)
)

func looksLikeUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

func looksLikeHex(s string) bool {
	return hexRegex.MatchString(strings.ToLower(s))
}

func looksLikeBase64(s string) bool {
	return base64Regex.MatchString(s)
}
