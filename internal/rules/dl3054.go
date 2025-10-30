package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3054.hs.
func DL3054() rule.Rule {
	return rule.NewSimpleRule(
		"DL3054",
		rule.Warning,
		"Label `",
		checkDL3054,
	)
}

func checkDL3054(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3054.hs
	// See: hadolint/src/Hadolint/Rule/DL3054.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
