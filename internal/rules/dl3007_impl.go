package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3007 creates a rule for checking FROM tag is not "latest".
func DL3007() rule.Rule {
	return rule.NewSimpleRule(
		DL3007Meta.Code,
		DL3007Meta.Severity,
		DL3007Meta.Message,
		checkDL3007,
	)
}

func checkDL3007(instruction syntax.Instruction) bool {
	from, ok := instruction.(*syntax.From)
	if !ok {
		return true
	}

	// If both tag and digest present - always OK
	if from.Image.Tag != nil && from.Image.Digest != nil {
		return true
	}

	// If tag present, must not be "latest"
	if from.Image.Tag != nil {
		return *from.Image.Tag != "latest"
	}

	// No tag - OK (DL3006 will handle this case)
	return true
}
