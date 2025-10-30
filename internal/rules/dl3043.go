package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3043.hs.
func DL3043() rule.Rule {
	return rule.NewSimpleRule(
		"DL3043",
		rule.Error,
		"`ONBUILD`, `FROM` or `MAINTAINER` triggered from within `ONBUILD` instruction.",
		checkDL3043,
	)
}

func checkDL3043(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3043.hs
	// See: hadolint/src/Hadolint/Rule/DL3043.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
