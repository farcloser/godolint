package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3025.hs.
func DL3025() rule.Rule {
	return rule.NewSimpleRule(
		"DL3025",
		rule.Warning,
		"Use arguments JSON notation for CMD and ENTRYPOINT arguments",
		checkDL3025,
	)
}

func checkDL3025(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3025.hs
	// See: hadolint/src/Hadolint/Rule/DL3025.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
