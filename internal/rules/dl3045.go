package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3045.hs.
func DL3045() rule.Rule {
	return rule.NewSimpleRule(
		"DL3045",
		rule.Warning,
		"`COPY` to a relative destination without `WORKDIR` set.",
		checkDL3045,
	)
}

func checkDL3045(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3045.hs
	// See: hadolint/src/Hadolint/Rule/DL3045.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
