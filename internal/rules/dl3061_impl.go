package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3061Rule checks that Dockerfile begins with FROM, ARG, or comment.
// Ported from Hadolint.Rule.DL3061.
type DL3061Rule struct{}

// DL3061 creates the rule for checking instruction order.
func DL3061() rule.Rule {
	return &DL3061Rule{}
}

func (r *DL3061Rule) Code() rule.RuleCode {
	return DL3061Meta.Code
}

func (r *DL3061Rule) Severity() rule.Severity {
	return DL3061Meta.Severity
}

func (r *DL3061Rule) Message() string {
	return DL3061Meta.Message
}

func (r *DL3061Rule) InitialState() rule.State {
	return rule.EmptyState(false) // false = haven't seen FROM yet
}

// Check ensures Dockerfile starts with FROM, ARG, or comment.
// Ported from the check function in DL3061.hs.
func (r *DL3061Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	seenFrom := state.Data.(bool)

	// Once we've seen FROM, everything is OK
	if seenFrom {
		return state
	}

	// FROM - mark that we've seen it
	if _, ok := instruction.(*syntax.From); ok {
		return state.ReplaceData(true)
	}

	// ARG before FROM - OK
	if _, ok := instruction.(*syntax.Arg); ok {
		return state
	}

	// Comments and pragmas before FROM - OK
	// Note: Our parser may not preserve comments as instructions
	// If comments aren't in the instruction stream, this rule will work correctly
	// because instructions before FROM that aren't ARG will fail

	// Any other instruction before FROM - fail
	return state.AddFailure(rule.CheckFailure{
		Code:     DL3061Meta.Code,
		Severity: DL3061Meta.Severity,
		Message:  DL3061Meta.Message,
		Line:     line,
	})
}

func (r *DL3061Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
