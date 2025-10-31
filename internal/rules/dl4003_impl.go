package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

type cmdState int

const (
	noCmd cmdState = iota
	hasCmd
)

// DL4003Rule is a stateful rule that tracks CMD instructions.
// Ported from Hadolint.Rule.DL4003.
type DL4003Rule struct{}

// DL4003 creates the rule for checking multiple CMD instructions.
func DL4003() rule.Rule {
	return &DL4003Rule{}
}

func (r *DL4003Rule) Code() rule.RuleCode {
	return DL4003Meta.Code
}

func (r *DL4003Rule) Severity() rule.Severity {
	return DL4003Meta.Severity
}

func (r *DL4003Rule) Message() string {
	return DL4003Meta.Message
}

func (r *DL4003Rule) InitialState() rule.State {
	return rule.EmptyState(noCmd)
}

// Check implements the stateful check for multiple CMD instructions.
// Ported from the check function in DL4003.hs.
func (r *DL4003Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	// Reset state on each FROM (new stage)
	if _, ok := instruction.(*syntax.From); ok {
		return state.ReplaceData(noCmd)
	}

	// Check for CMD instruction
	if _, ok := instruction.(*syntax.Cmd); ok {
		currentState := state.Data.(cmdState)

		if currentState == noCmd {
			// First CMD in this stage - remember it
			return state.ReplaceData(hasCmd)
		}

		// Second CMD in this stage - add failure
		return state.AddFailure(rule.CheckFailure{
			Code:     DL4003Meta.Code,
			Severity: DL4003Meta.Severity,
			Message:  DL4003Meta.Message,
			Line:     line,
		})
	}

	return state
}

func (r *DL4003Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
