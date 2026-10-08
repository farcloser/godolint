package rules

import (
	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/syntax"
)

// DL3065 creates the rule that FROM --platform is not the default
// $TARGETPLATFORM. Ported from Hadolint.Rule.DL3065.
func DL3065() rule.Rule {
	return rule.NewSimpleRule(
		DL3065Meta.Code,
		DL3065Meta.Severity,
		DL3065Meta.Message,
		checkDL3065,
	)
}

// checkDL3065 fails a FROM whose --platform names the predefined
// $TARGETPLATFORM, in either spelling, since that is what BuildKit does
// without the flag.
func checkDL3065(instruction syntax.Instruction) bool {
	from, ok := instruction.(*syntax.From)
	if !ok || from.Image.Platform == nil {
		return true
	}

	platform := *from.Image.Platform

	return platform != "$TARGETPLATFORM" && platform != "${TARGETPLATFORM}"
}
