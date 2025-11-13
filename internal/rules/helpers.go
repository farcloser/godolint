package rules

import (
	"strings"
	"unicode"
)

// dropQuotes removes surrounding quotes from a string.
func dropQuotes(str string) string {
	str = strings.TrimSpace(str)
	if len(str) >= 2 {
		if (str[0] == '"' && str[len(str)-1] == '"') || (str[0] == '\'' && str[len(str)-1] == '\'') {
			return str[1 : len(str)-1]
		}
	}

	return str
}

// isWindowsAbsolute checks if a path is a Windows absolute path (e.g., C:\).
func isWindowsAbsolute(path string) bool {
	if len(path) < 2 {
		return false
	}

	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}
