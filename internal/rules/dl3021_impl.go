package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3021 creates a rule for checking COPY with multiple sources ends with /.
func DL3021() rule.Rule {
	return rule.NewSimpleRule(
		DL3021Meta.Code,
		DL3021Meta.Severity,
		DL3021Meta.Message,
		checkDL3021,
	)
}

func checkDL3021(instruction syntax.Instruction) bool {
	copyInstr, ok := instruction.(*syntax.Copy)
	if !ok {
		return true
	}

	// If only one source, no requirement
	if len(copyInstr.Source) <= 1 {
		return true
	}

	// Multiple sources - destination must end with /
	dest := dropQuotes(copyInstr.Destination)
	if len(dest) == 0 {
		return false
	}

	return dest[len(dest)-1] == '/'
}
