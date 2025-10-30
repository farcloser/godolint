package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL4004.hs.
func DL4004() rule.Rule {
	return rule.NewSimpleRule(
		"DL4004",
		rule.Error,
		"Multiple `ENTRYPOINT` instructions found. If you list more than one `ENTRYPOINT` then       only the last `ENTRYPOINT` will take effect",
		checkDL4004,
	)
}

func checkDL4004(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL4004.hs
	// See: hadolint/src/Hadolint/Rule/DL4004.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
