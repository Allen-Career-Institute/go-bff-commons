package utils

import (
	"strings"
)

func Join(strs ...string) string {
	var sb strings.Builder

	for _, str := range strs {
		sb.WriteString(str)
	}

	return sb.String()
}

// JoinRunes you can use the WriteRune and WriteByte methods to add single characters to your string as you build it.
func JoinRunes(runes ...rune) string {
	var sb strings.Builder

	for _, r := range runes {
		sb.WriteRune(r)
	}

	return sb.String()
}
