package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3045 checks for COPY to relative destination without WORKDIR set.
func DL3045() rule.Rule {
	return &DL3045Rule{}
}

// DL3045Rule implements the rule for checking COPY with relative paths.
type DL3045Rule struct{}

// Code returns the rule code.
func (*DL3045Rule) Code() rule.RuleCode {
	return DL3045Meta.Code
}

// Severity returns the rule severity.
func (*DL3045Rule) Severity() rule.Severity {
	return DL3045Meta.Severity
}

// Message returns the rule message.
func (*DL3045Rule) Message() string {
	return DL3045Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3045Rule) InitialState() rule.State {
	return rule.EmptyState(dl3045State{
		WorkdirSet: make(map[string]bool),
	})
}

type dl3045State struct {
	CurrentStage string
	// Map from stage name to whether WORKDIR has been set
	WorkdirSet map[string]bool
}

// Check validates COPY instructions use absolute paths or have WORKDIR set.
func (r *DL3045Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	var currentState dl3045State
	if state.Data != nil {
		currentState = state.Data.(dl3045State)
	} else {
		currentState = dl3045State{
			WorkdirSet: make(map[string]bool),
		}
	}

	switch inst := instruction.(type) {
	case *syntax.From:
		// Remember the stage
		stageName := inst.Image.Image
		if inst.Image.Alias != nil {
			stageName = *inst.Image.Alias
		}

		// Check if we inherit WORKDIR from parent stage
		inheritWorkdir := false
		if currentState.WorkdirSet[inst.Image.Image] {
			inheritWorkdir = true
		}

		currentState.CurrentStage = stageName
		currentState.WorkdirSet[stageName] = inheritWorkdir

		return state.ReplaceData(currentState)

	case *syntax.Workdir:
		// Mark that WORKDIR has been set in current stage
		currentState.WorkdirSet[currentState.CurrentStage] = true

		return state.ReplaceData(currentState)

	case *syntax.Copy:
		// Get destination
		dest := inst.Destination
		dest = strings.Trim(dest, "\"'")

		// Don't fail if:
		// 1. WORKDIR has been set
		if currentState.WorkdirSet[currentState.CurrentStage] {
			return state.ReplaceData(currentState)
		}

		// 2. Destination is absolute (starts with /)
		if strings.HasPrefix(dest, "/") {
			return state.ReplaceData(currentState)
		}

		// 3. Destination is Windows absolute (like C:\path)
		if len(dest) >= 2 && dest[1] == ':' && isLetter(rune(dest[0])) {
			return state.ReplaceData(currentState)
		}

		// 4. Destination is a variable (starts with $)
		if strings.HasPrefix(dest, "$") {
			return state.ReplaceData(currentState)
		}

		// Otherwise, this is a violation
		return state.
			ReplaceData(currentState).
			AddFailure(rule.CheckFailure{
				Code:     r.Code(),
				Severity: r.Severity(),
				Message:  r.Message(),
				Line:     line,
			})
	}

	return state.ReplaceData(currentState)
}

// Finalize performs final checks after processing all instructions.
func (*DL3045Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
