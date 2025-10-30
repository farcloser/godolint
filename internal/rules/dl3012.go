package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3012.hs.
func DL3012() rule.Rule {
	return rule.NewSimpleRule(
		"DL3012",
		rule.Error,
		"Multiple `HEALTHCHECK` instructions",
		checkDL3012,
	)
}

func checkDL3012(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3012.hs
	// See: hadolint/src/Hadolint/Rule/DL3012.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
