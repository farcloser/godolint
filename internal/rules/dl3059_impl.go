package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3059 creates a rule for checking multiple consecutive RUN instructions.
func DL3059() rule.Rule {
	return &DL3059Rule{}
}

// DL3059Rule checks for multiple consecutive RUN instructions.
// Ported from Hadolint.Rule.DL3059.
type DL3059Rule struct{}

// dl3059State tracks the previous RUN instruction state.
type dl3059State struct {
	Flags []string // Flags from previous RUN
	Count int      // Number of commands in previous RUN
}

func (r *DL3059Rule) Code() rule.RuleCode {
	return DL3059Meta.Code
}

func (r *DL3059Rule) Severity() rule.Severity {
	return DL3059Meta.Severity
}

func (r *DL3059Rule) Message() string {
	return DL3059Meta.Message
}

func (r *DL3059Rule) InitialState() rule.State {
	return rule.EmptyState(dl3059State{})
}

func (r *DL3059Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	var prevState dl3059State
	if state.Data != nil {
		prevState = state.Data.(dl3059State)
	}

	switch instr := instruction.(type) {
	case *syntax.Run:
		// Parse the shell command to count commands
		parsed, err := shell.ParseShell(instr.Command)
		if err != nil {
			// If parsing fails, reset state and continue
			return state.ReplaceData(dl3059State{})
		}

		commandCount := shell.CountCommands(parsed)
		currentFlags := instr.Flags

		// Check if we should fail:
		// 1. Previous state exists (not first RUN)
		// 2. Flags are the same
		// 3. NEITHER current NOR previous have chained commands (2+)
		// Per hadolint: "ignores the case where multiple commands are chained together"
		if prevState.Count > 0 && // Previous RUN exists
			flagsEqual(prevState.Flags, currentFlags) &&
			commandCount < 2 && // Current is not chained
			prevState.Count < 2 { // Previous is not chained
			// Fail: consecutive simple RUN instructions with same flags
			return state.AddFailure(rule.CheckFailure{
				Code:     r.Code(),
				Severity: r.Severity(),
				Message:  r.Message(),
				Line:     line,
			}).ReplaceData(dl3059State{
				Flags: currentFlags,
				Count: commandCount,
			})
		}

		// Update state with current RUN
		return state.ReplaceData(dl3059State{
			Flags: currentFlags,
			Count: commandCount,
		})

	case *syntax.Comment:
		// Comments don't break consecutive RUN detection
		return state

	default:
		// Any other instruction resets the state
		return state.ReplaceData(dl3059State{})
	}
}

func (r *DL3059Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}

// flagsEqual compares two flag slices for equality.
func flagsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Simple comparison - flags should be in same order
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
