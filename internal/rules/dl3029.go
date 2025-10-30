package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3029.hs.
func DL3029() rule.Rule {
	return rule.NewSimpleRule(
		"DL3029",
		rule.Warning,
		"Do not use --platform flag with FROM",
		checkDL3029,
	)
}

func checkDL3029(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3029.hs
	// See: hadolint/src/Hadolint/Rule/DL3029.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
