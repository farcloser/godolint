package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3011.hs.
func DL3011() rule.Rule {
	return rule.NewSimpleRule(
		"DL3011",
		rule.Error,
		"Valid UNIX ports range from 0 to 65535",
		checkDL3011,
	)
}

func checkDL3011(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3011.hs
	// See: hadolint/src/Hadolint/Rule/DL3011.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
