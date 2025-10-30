package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3055.hs.
func DL3055() rule.Rule {
	return rule.NewSimpleRule(
		"DL3055",
		rule.Warning,
		"Label `",
		checkDL3055,
	)
}

func checkDL3055(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3055.hs
	// See: hadolint/src/Hadolint/Rule/DL3055.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
