package rules

import (
	"regexp"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL1001 checks for inline ignore pragmas.
func DL1001() rule.Rule {
	return rule.NewSimpleRule(
		DL1001Meta.Code,
		DL1001Meta.Severity,
		DL1001Meta.Message,
		checkDL1001,
	)
}

func checkDL1001(instruction syntax.Instruction) bool {
	comment, ok := instruction.(*syntax.Comment)
	if !ok {
		return true // Not a comment, pass
	}

	// Check if comment contains ignore pragma
	return !isIgnorePragma(comment.Text)
}

// isIgnorePragma checks if a comment text contains a hadolint ignore pragma.
// Matches patterns like:
// - hadolint ignore=DL3000
// - hadolint ignore=DL3000,DL3001
func isIgnorePragma(text string) bool {
	// Match hadolint ignore= followed by rule codes
	pattern := regexp.MustCompile(`hadolint\s+ignore\s*=\s*DL\d{4}`)
	return pattern.MatchString(text)
}
