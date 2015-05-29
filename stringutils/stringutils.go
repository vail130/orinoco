package stringutils

import (
	"bytes"
	"strings"
)

func Concat(strings ...string) string {
	var buffer bytes.Buffer

	for _, str := range strings {
		buffer.WriteString(str)
	}

	return buffer.String()
}

func UnderscoreToTitle(s string) string {
	stringParts := strings.Split(s, "_")
	for i := 0; i < len(stringParts); i++ {
		stringParts[i] = strings.Title(stringParts[i])
	}
	return strings.Join(stringParts, "")
}

func StringToBool(s string) bool {
	s = strings.ToLower(s)
	if s == "" || s == "f" || s == "false" || s == "0" {
		return false
	}

	return true
}
