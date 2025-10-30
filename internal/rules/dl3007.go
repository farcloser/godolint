// Package rules implements individual hadolint rules.
package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Ported from Hadolint.Rule.DL3007.
func DL3007() rule.Rule {
	return rule.NewSimpleRule(
		"DL3007",
		rule.Warning,
		"Using latest is prone to errors if the image will ever update. Pin the version explicitly to a release tag",
		checkDL3007,
	)
}

func checkDL3007(instruction syntax.Instruction) bool {
	from, ok := instruction.(*syntax.From)
	if !ok {
		return true
	}

	// If both tag and digest are present, digest overrides tag -> pass
	if from.Image.Tag != nil && from.Image.Digest != nil {
		return true
	}

	// If tag is "latest" (without digest) -> fail
	if from.Image.Tag != nil && *from.Image.Tag == "latest" {
		return false
	}

	// Otherwise pass
	return true
}
