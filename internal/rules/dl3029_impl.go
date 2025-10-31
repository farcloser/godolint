package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3029 creates a rule for checking --platform flag usage.
func DL3029() rule.Rule {
	return rule.NewSimpleRule(
		DL3029Meta.Code,
		DL3029Meta.Severity,
		DL3029Meta.Message,
		checkDL3029,
	)
}

func checkDL3029(instruction syntax.Instruction) bool {
	from, ok := instruction.(*syntax.From)
	if !ok {
		return true
	}

	// No platform flag - OK
	if from.Image.Platform == nil {
		return true
	}

	platform := *from.Image.Platform

	// Platform contains BUILDPLATFORM or TARGETPLATFORM - OK (build args)
	if strings.Contains(platform, "BUILDPLATFORM") || strings.Contains(platform, "TARGETPLATFORM") {
		return true
	}

	// Platform flag with hardcoded value - fail
	return false
}
