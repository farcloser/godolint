package pragma

import "regexp"

// shellPragma is hadolint's shell pragma, "hadolint shell=powershell": spaces
// or tabs around the keywords and the "=", a shell name running to the next
// space, and an optional "# comment" before the end. Ported from shellParser
// in Hadolint.Pragma.
var shellPragma = regexp.MustCompile(`^[ \t]*hadolint[ \t]+shell[ \t]*=[ \t]*([^ \t\n]+)[ \t]*(?:#[^\n]*)?$`)

// ParseShell reads a shell pragma and returns the shell it names.
func ParseShell(text string) (string, bool) {
	match := shellPragma.FindStringSubmatch(text)
	if match == nil {
		return "", false
	}

	return match[1], true
}
