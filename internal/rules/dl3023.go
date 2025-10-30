package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3023.hs.
func DL3023() rule.Rule {
	return rule.NewSimpleRule(
		"DL3023",
		rule.Error,
		"`COPY --from` cannot reference its own `FROM` alias",
		checkDL3023,
	)
}

func checkDL3023(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/DL3023.hs
	// See: hadolint/src/Hadolint/Rule/DL3023.hs for implementation
	// Placeholder: allow all instructions for now
	return true
}
