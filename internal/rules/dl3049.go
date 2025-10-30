package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3049.hs.
func DL3049() rule.Rule {
	return rule.NewSimpleRule(
		"DL3049",
		rule.Info,
		"Label `",
		checkDL3049,
	)
}

func checkDL3049(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3049.hs
	// See: hadolint/src/Hadolint/Rule/DL3049.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
