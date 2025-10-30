package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Ported from Hadolint.Rule.DL4000.
func DL4000() rule.Rule {
	return rule.NewSimpleRule(
		"DL4000",
		rule.Error,
		"MAINTAINER is deprecated",
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
