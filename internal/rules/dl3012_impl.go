package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

type healthcheckState int

const (
	noHealthcheck healthcheckState = iota
	hasHealthcheck
)

// DL3012Rule is a stateful rule that tracks HEALTHCHECK instructions.
// Ported from Hadolint.Rule.DL3012.
type DL3012Rule struct{}

// DL3012 creates the rule for checking multiple HEALTHCHECK instructions.
func DL3012() rule.Rule {
	return &DL3012Rule{}
}

// Code returns the rule code.
func (*DL3012Rule) Code() rule.RuleCode {
	return DL3012Meta.Code
}

// Severity returns the rule severity.
func (*DL3012Rule) Severity() rule.Severity {
	return DL3012Meta.Severity
}

// Message returns the rule message.
func (*DL3012Rule) Message() string {
	return DL3012Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3012Rule) InitialState() rule.State {
	return rule.EmptyState(noHealthcheck)
}

// Check implements the stateful check for multiple HEALTHCHECK instructions.
// Ported from the check function in DL3012.hs.
func (*DL3012Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	// Reset state on each FROM (new stage)
	if _, ok := instruction.(*syntax.From); ok {
		return state.ReplaceData(noHealthcheck)
	}

	// Check for HEALTHCHECK instruction
	if _, ok := instruction.(*syntax.Healthcheck); ok {
		currentState := state.Data.(healthcheckState)

		if currentState == noHealthcheck {
			// First HEALTHCHECK in this stage - remember it
			return state.ReplaceData(hasHealthcheck)
		}

		// Second HEALTHCHECK in this stage - add failure
		return state.AddFailure(rule.CheckFailure{
			Code:     DL3012Meta.Code,
			Severity: DL3012Meta.Severity,
			Message:  DL3012Meta.Message,
			Line:     line,
			Column:   1, // Hardcoded to 1 (matches hadolint)
		})
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*DL3012Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
