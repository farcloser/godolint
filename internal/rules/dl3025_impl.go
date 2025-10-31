package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3025 creates a rule for checking CMD and ENTRYPOINT use JSON notation.
func DL3025() rule.Rule {
	return rule.NewSimpleRule(
		DL3025Meta.Code,
		DL3025Meta.Severity,
		DL3025Meta.Message,
		checkDL3025,
	)
}

func checkDL3025(instruction syntax.Instruction) bool {
	switch inst := instruction.(type) {
	case *syntax.Cmd:
		// Fail if using shell form (not JSON)
		return inst.IsJSON
	case *syntax.Entrypoint:
		// Fail if using shell form (not JSON)
		return inst.IsJSON
	}

	return true
}
