package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3053.hs.
func DL3053() rule.Rule {
	return rule.NewSimpleRule(
		"DL3053",
		rule.Warning,
		"Label `",
		checkDL3053,
	)
}

func checkDL3053(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3053.hs
	// See: hadolint/src/Hadolint/Rule/DL3053.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
