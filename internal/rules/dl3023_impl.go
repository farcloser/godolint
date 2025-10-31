package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3023Rule checks that COPY --from doesn't reference its own FROM alias.
// Ported from Hadolint.Rule.DL3023.
type DL3023Rule struct{}

// DL3023 creates the rule for checking COPY --from self-reference.
func DL3023() rule.Rule {
	return &DL3023Rule{}
}

func (r *DL3023Rule) Code() rule.RuleCode {
	return DL3023Meta.Code
}

func (r *DL3023Rule) Severity() rule.Severity {
	return DL3023Meta.Severity
}

func (r *DL3023Rule) Message() string {
	return DL3023Meta.Message
}

func (r *DL3023Rule) InitialState() rule.State {
	return rule.EmptyState("") // Empty string = no current stage alias
}

// Check tracks current stage alias and validates COPY --from doesn't self-reference.
// Ported from the check function in DL3023.hs.
func (r *DL3023Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	currentAlias := state.Data.(string)

	// Remember current FROM alias
	if from, ok := instruction.(*syntax.From); ok {
		if from.Image.Alias != nil {
			return state.ReplaceData(*from.Image.Alias)
		}

		return state.ReplaceData("") // No alias
	}

	// Check COPY --from doesn't reference current stage
	if copy, ok := instruction.(*syntax.Copy); ok {
		if copy.From != nil && currentAlias != "" {
			if *copy.From == currentAlias {
				// Self-reference - fail
				return state.AddFailure(rule.CheckFailure{
					Code:     DL3023Meta.Code,
					Severity: DL3023Meta.Severity,
					Message:  DL3023Meta.Message,
					Line:     line,
				})
			}
		}
	}

	return state
}

func (r *DL3023Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
