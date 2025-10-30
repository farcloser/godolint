package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3021.hs.
func DL3021() rule.Rule {
	return rule.NewSimpleRule(
		"DL3021",
		rule.Error,
		"COPY with more than 2 arguments requires the last argument to end with /",
		checkDL3021,
	)
}

func checkDL3021(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3021.hs
	// See: hadolint/src/Hadolint/Rule/DL3021.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
