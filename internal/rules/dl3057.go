package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3057.hs.
func DL3057() rule.Rule {
	return rule.NewSimpleRule(
		"DL3057",
		rule.Warning,
		"`HEALTHCHECK` instruction missing.",
		checkDL3057,
	)
}

func checkDL3057(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3057.hs
	// See: hadolint/src/Hadolint/Rule/DL3057.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
