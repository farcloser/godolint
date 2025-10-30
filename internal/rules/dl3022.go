package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3022.hs.
func DL3022() rule.Rule {
	return rule.NewSimpleRule(
		"DL3022",
		rule.Warning,
		"`COPY --from` should reference a previously defined `FROM` alias",
		checkDL3022,
	)
}

func checkDL3022(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3022.hs
	// See: hadolint/src/Hadolint/Rule/DL3022.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
