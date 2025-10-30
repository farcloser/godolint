package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3024.hs.
func DL3024() rule.Rule {
	return rule.NewSimpleRule(
		"DL3024",
		rule.Error,
		"FROM aliases (stage names) must be unique",
		checkDL3024,
	)
}

func checkDL3024(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3024.hs
	// See: hadolint/src/Hadolint/Rule/DL3024.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
