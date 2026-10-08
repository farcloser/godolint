package rules

import (
	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/syntax"
)

// DL3067 creates the rule that a stage's whole filesystem is not copied.
// Ported from Hadolint.Rule.DL3067.
func DL3067() rule.Rule {
	return rule.NewSimpleRule(
		DL3067Meta.Code,
		DL3067Meta.Severity,
		DL3067Meta.Message,
		checkDL3067,
	)
}

// checkDL3067 fails a COPY --from whose sources include the root and whose
// destination is the root: the other stage's filesystem copied whole.
func checkDL3067(instruction syntax.Instruction) bool {
	copyInstr, ok := instruction.(*syntax.Copy)
	if !ok || copyInstr.From == nil {
		return true
	}

	rootSource := false

	for _, source := range copyInstr.Source {
		if dropQuotes(source) == "/" {
			rootSource = true

			break
		}
	}

	return !(rootSource && dropQuotes(copyInstr.Destination) == "/")
}
