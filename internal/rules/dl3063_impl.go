package rules

import (
	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/syntax"
)

// DL3063 creates the rule that a stage name is not a reserved word.
// Ported from Hadolint.Rule.DL3063.
func DL3063() rule.Rule {
	return rule.NewSimpleRule(
		DL3063Meta.Code,
		DL3063Meta.Severity,
		DL3063Meta.Message,
		checkDL3063,
	)
}

// checkDL3063 fails a FROM whose alias is a name BuildKit reserves: `scratch`
// is the empty image, `context` the build context.
func checkDL3063(instruction syntax.Instruction) bool {
	from, ok := instruction.(*syntax.From)
	if !ok || from.Image.Alias == nil {
		return true
	}

	switch *from.Image.Alias {
	case scratchImage, "context":
		return false
	default:
		return true
	}
}
