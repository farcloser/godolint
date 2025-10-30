package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3050.hs.
func DL3050() rule.Rule {
	return rule.NewSimpleRule(
		"DL3050",
		rule.Info,
		"Superfluous label(s) present.",
		checkDL3050,
	)
}

func checkDL3050(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3050.hs
	// See: hadolint/src/Hadolint/Rule/DL3050.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
