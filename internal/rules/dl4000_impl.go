package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL4000 creates a rule from the generated metadata.
// Auto-generated from hadolint Haskell source.
func DL4000() rule.Rule {
	return rule.NewSimpleRule(
		DL4000Meta.Code,
		DL4000Meta.Severity,
		DL4000Meta.Message,
		checkDL4000,
	)
}

func checkDL4000(instruction syntax.Instruction) bool {
	_, ok := instruction.(*syntax.Maintainer)
	if ok {
		return false // Maintainer instruction found -> fail
	}

	return true
}
