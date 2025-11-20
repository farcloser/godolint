package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

type entrypointState int

const (
	noEntrypoint entrypointState = iota
	hasEntrypoint
)

// DL4004Rule is a stateful rule that tracks ENTRYPOINT instructions.
// Ported from Hadolint.Rule.DL4004.
type DL4004Rule struct{}

// DL4004 creates the rule for checking multiple ENTRYPOINT instructions.
func DL4004() rule.Rule {
	return &DL4004Rule{}
}

// Code returns the rule code.
func (*DL4004Rule) Code() rule.RuleCode {
	return DL4004Meta.Code
}

// Severity returns the rule severity.
func (*DL4004Rule) Severity() rule.Severity {
	return DL4004Meta.Severity
}

// Message returns the rule message.
func (*DL4004Rule) Message() string {
	return DL4004Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL4004Rule) InitialState() rule.State {
	return rule.EmptyState(noEntrypoint)
}

// Check implements the stateful check for multiple ENTRYPOINT instructions.
// Ported from the check function in DL4004.hs.
func (*DL4004Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	// Reset state on each FROM (new stage)
	if _, ok := instruction.(*syntax.From); ok {
		return state.ReplaceData(noEntrypoint)
	}

	// Check for ENTRYPOINT instruction
	if _, ok := instruction.(*syntax.Entrypoint); ok {
		currentState := state.Data.(entrypointState)

		if currentState == noEntrypoint {
			// First ENTRYPOINT in this stage - remember it
			return state.ReplaceData(hasEntrypoint)
		}

		// Second ENTRYPOINT in this stage - add failure
		return state.AddFailure(rule.CheckFailure{
			Code:     DL4004Meta.Code,
			Severity: DL4004Meta.Severity,
			Message:  DL4004Meta.Message,
			Line:     line,
			Column:   1, // Hardcoded to 1 (matches hadolint)
		})
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*DL4004Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
