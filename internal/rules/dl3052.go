package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3052.hs.
func DL3052() rule.Rule {
	return rule.NewSimpleRule(
		"DL3052",
		rule.Warning,
		"Label `",
		checkDL3052,
	)
}

func checkDL3052(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3052.hs
	// See: hadolint/src/Hadolint/Rule/DL3052.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
