package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL4001 checks for using both wget and curl.
func DL4001() rule.Rule {
	return &DL4001Rule{}
}

// DL4001Rule implements the rule for detecting mixed usage of wget and curl.
type DL4001Rule struct{}

// Code returns the rule code.
func (*DL4001Rule) Code() rule.RuleCode {
	return DL4001Meta.Code
}

// Severity returns the rule severity.
func (*DL4001Rule) Severity() rule.Severity {
	return DL4001Meta.Severity
}

// Message returns the rule message.
func (*DL4001Rule) Message() string {
	return DL4001Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL4001Rule) InitialState() rule.State {
	return rule.EmptyState(dl4001State{})
}

type dl4001State struct {
	HasCurl bool
	HasWget bool
}

// Check detects if both wget and curl are used in the same Dockerfile.
func (r *DL4001Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	var currentState dl4001State
	if state.Data != nil {
		currentState = state.Data.(dl4001State)
	}

	switch inst := instruction.(type) {
	case *syntax.From:
		// Reset state for each stage
		return state.ReplaceData(dl4001State{})

	case *syntax.Run:
		parsed, err := shell.ParseShell(inst.Command)
		if err != nil {
			return state.ReplaceData(currentState)
		}

		commands := shell.FindCommandNames(parsed)
		newCurl := currentState.HasCurl
		newWget := currentState.HasWget

		for _, cmd := range commands {
			switch cmd {
			case "curl":
				newCurl = true
			case "wget":
				newWget = true
			}
		}

		// If we just found both in this RUN or now have both, fail
		if newCurl && newWget && (!currentState.HasCurl || !currentState.HasWget) {
			return state.
				ReplaceData(dl4001State{HasCurl: newCurl, HasWget: newWget}).
				AddFailure(rule.CheckFailure{
					Code:     r.Code(),
					Severity: r.Severity(),
					Message:  r.Message(),
					Line:     line,
				})
		}

		return state.ReplaceData(dl4001State{HasCurl: newCurl, HasWget: newWget})
	}

	return state.ReplaceData(currentState)
}

// Finalize performs final checks after processing all instructions.
func (*DL4001Rule) Finalize(state rule.State) rule.State {
	return state
}
