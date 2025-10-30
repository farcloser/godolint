package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3026.hs.
func DL3026() rule.Rule {
	return rule.NewSimpleRule(
		"DL3026",
		rule.Error,
		"Use only an allowed registry in the FROM image",
		checkDL3026,
	)
}

func checkDL3026(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3026.hs
	// See: hadolint/src/Hadolint/Rule/DL3026.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
