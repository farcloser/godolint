package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3051.hs.
func DL3051() rule.Rule {
	return rule.NewSimpleRule(
		"DL3051",
		rule.Warning,
		"label `",
		checkDL3051,
	)
}

func checkDL3051(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3051.hs
	// See: hadolint/src/Hadolint/Rule/DL3051.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
