package rules

import (
	"strings"
	"unicode"
)

// dropQuotes removes surrounding quotes from a string.
func dropQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}

	return s
}

// isWindowsAbsolute checks if a path is a Windows absolute path (e.g., C:\).
func isWindowsAbsolute(path string) bool {
	if len(path) < 2 {
		return false
	}

	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}
