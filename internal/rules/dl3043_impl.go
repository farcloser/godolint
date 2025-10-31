package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3043 creates a rule from the generated metadata.
// Checks for invalid ONBUILD instructions.
func DL3043() rule.Rule {
	return rule.NewSimpleRule(
		DL3043Meta.Code,
		DL3043Meta.Severity,
		DL3043Meta.Message,
		checkDL3043,
	)
}

func checkDL3043(instruction syntax.Instruction) bool {
	onbuild, ok := instruction.(*syntax.OnBuild)
	if !ok {
		return true
	}

	if onbuild.Inner == nil {
		return true
	}

	// Check for ONBUILD ONBUILD
	if _, ok := onbuild.Inner.(*syntax.OnBuild); ok {
		return false
	}

	// Check for ONBUILD FROM
	if _, ok := onbuild.Inner.(*syntax.From); ok {
		return false
	}

	// Check for ONBUILD MAINTAINER
	if _, ok := onbuild.Inner.(*syntax.Maintainer); ok {
		return false
	}

	return true
}
