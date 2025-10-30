package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3058.hs.
func DL3058() rule.Rule {
	return rule.NewSimpleRule(
		"DL3058",
		rule.Warning,
		"Label `",
		checkDL3058,
	)
}

func checkDL3058(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3058.hs
	// See: hadolint/src/Hadolint/Rule/DL3058.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
