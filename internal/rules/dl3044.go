package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3044.hs.
func DL3044() rule.Rule {
	return rule.NewSimpleRule(
		"DL3044",
		rule.Error,
		"Do not refer to an environment variable within the same `ENV` statement where it is defined.",
		checkDL3044,
	)
}

func checkDL3044(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3044.hs
	// See: hadolint/src/Hadolint/Rule/DL3044.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
