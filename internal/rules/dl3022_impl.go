package rules

import (
	"strconv"
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3022State tracks FROM stages to validate COPY --from references.
type dl3022State struct {
	count  int             // Number of FROM stages seen
	stages map[string]bool // Named stages (aliases)
}

// DL3022Rule checks that COPY --from references valid stages.
// Ported from Hadolint.Rule.DL3022.
type DL3022Rule struct{}

// DL3022 creates the rule for checking COPY --from references.
func DL3022() rule.Rule {
	return &DL3022Rule{}
}

func (r *DL3022Rule) Code() rule.RuleCode {
	return DL3022Meta.Code
}

func (r *DL3022Rule) Severity() rule.Severity {
	return DL3022Meta.Severity
}

func (r *DL3022Rule) Message() string {
	return DL3022Meta.Message
}

func (r *DL3022Rule) InitialState() rule.State {
	return rule.EmptyState(dl3022State{
		count:  0,
		stages: make(map[string]bool),
	})
}

// Check tracks FROM stages and validates COPY --from references.
// Ported from the check function in DL3022.hs.
func (r *DL3022Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3022State)

	// Track FROM stages
	if from, ok := instruction.(*syntax.From); ok {
		newStages := make(map[string]bool)
		for k, v := range s.stages {
			newStages[k] = v
		}

		if from.Image.Alias != nil {
			newStages[*from.Image.Alias] = true
		}

		s.count++
		s.stages = newStages

		return state.ReplaceData(s)
	}

	// Check COPY --from references
	if copy, ok := instruction.(*syntax.Copy); ok {
		if copy.From == nil {
			return state
		}

		fromRef := *copy.From

		// Image reference (contains :) - OK (external image)
		if strings.Contains(fromRef, ":") {
			return state
		}

		// Named stage - OK if exists
		if s.stages[fromRef] {
			return state
		}

		// Numeric stage reference - OK if valid
		if idx, err := strconv.Atoi(fromRef); err == nil {
			if idx < s.count {
				return state
			}
		}

		// Invalid reference - fail
		return state.AddFailure(rule.CheckFailure{
			Code:     DL3022Meta.Code,
			Severity: DL3022Meta.Severity,
			Message:  DL3022Meta.Message,
			Line:     line,
		})
	}

	return state
}

func (r *DL3022Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
