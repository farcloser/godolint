package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Ported from Hadolint.Rule.DL3000.
func DL3000() rule.Rule {
	return rule.NewSimpleRule(
		"DL3000",
		rule.Error,
		"Use absolute WORKDIR",
		checkDL3000,
	)
}

func checkDL3000(instruction syntax.Instruction) bool {
	workdir, ok := instruction.(*syntax.Workdir)
	if !ok {
		return true
	}

	dir := dropQuotes(workdir.Directory)

	// Allow paths starting with environment variable
	if strings.HasPrefix(dir, "$") {
		return true
	}

	// Allow Unix absolute paths (/)
	if strings.HasPrefix(dir, "/") {
		return true
	}

	// Allow Windows absolute paths (C:\, D:\, etc.)
	if isWindowsAbsolute(dir) {
		return true
	}

	return false
}
