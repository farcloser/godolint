package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL4003.hs.
func DL4003() rule.Rule {
	return rule.NewSimpleRule(
		"DL4003",
		rule.Warning,
		"Multiple `CMD` instructions found. If you list more than one `CMD` then only the last       `CMD` will take effect",
		checkDL4003,
	)
}

func checkDL4003(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL4003.hs
	// See: hadolint/src/Hadolint/Rule/DL4003.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
