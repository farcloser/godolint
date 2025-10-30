package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3048.hs.
func DL3048() rule.Rule {
	return rule.NewSimpleRule(
		"DL3048",
		rule.Style,
		"Invalid label key.",
		checkDL3048,
	)
}

func checkDL3048(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3048.hs
	// See: hadolint/src/Hadolint/Rule/DL3048.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
