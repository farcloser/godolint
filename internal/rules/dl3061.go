package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3061.hs.
func DL3061() rule.Rule {
	return rule.NewSimpleRule(
		"DL3061",
		rule.Error,
		"Invalid instruction order. Dockerfile must begin with `FROM`,               `ARG` or comment.",
		checkDL3061,
	)
}

func checkDL3061(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3061.hs
	// See: hadolint/src/Hadolint/Rule/DL3061.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
