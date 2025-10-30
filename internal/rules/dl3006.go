package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3006.hs.
func DL3006() rule.Rule {
	return rule.NewSimpleRule(
		"DL3006",
		rule.Warning,
		"Always tag the version of an image explicitly",
		checkDL3006,
	)
}

func checkDL3006(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3006.hs
	// See: hadolint/src/Hadolint/Rule/DL3006.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
