package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3002.hs.
func DL3002() rule.Rule {
	return rule.NewSimpleRule(
		"DL3002",
		rule.Warning,
		"Last USER should not be root",
		checkDL3002,
	)
}

func checkDL3002(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3002.hs
	// See: hadolint/src/Hadolint/Rule/DL3002.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
