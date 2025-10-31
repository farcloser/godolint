package rules

import (
	"fmt"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3049State tracks which required labels have been defined.
type dl3049State struct {
	definedLabels map[string]bool // labels that have been defined
}

// DL3049Rule checks for missing required labels.
type DL3049Rule struct {
	cfg *config.Config
}

// DL3049 creates the rule for checking missing labels.
// TODO: Add full multi-stage tracking with inheritance
func DL3049() rule.Rule {
	return &DL3049Rule{
		cfg: config.Default(),
	}
}

func (r *DL3049Rule) Code() rule.RuleCode {
	return DL3049Meta.Code
}

func (r *DL3049Rule) Severity() rule.Severity {
	return DL3049Meta.Severity
}

func (r *DL3049Rule) Message() string {
	return DL3049Meta.Message
}

func (r *DL3049Rule) InitialState() rule.State {
	return rule.EmptyState(dl3049State{
		definedLabels: make(map[string]bool),
	})
}

func (r *DL3049Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3049State)

	// Track labels as they're defined
	if label, ok := instruction.(*syntax.Label); ok {
		for _, pair := range label.Pairs {
			s.definedLabels[pair.Key] = true
		}
		return state.ReplaceData(s)
	}

	// Reset on new FROM (simplified - full version would track inheritance)
	if _, ok := instruction.(*syntax.From); ok {
		s.definedLabels = make(map[string]bool)
		return state.ReplaceData(s)
	}

	return state
}

func (r *DL3049Rule) Finalize(state rule.State) rule.State {
	s := state.Data.(dl3049State)

	// Check for missing required labels
	for requiredLabel := range r.cfg.LabelSchema {
		if !s.definedLabels[requiredLabel] {
			// Report as missing (line 0 = end of file)
			state = state.AddFailure(rule.CheckFailure{
				Code:     DL3049Meta.Code,
				Severity: DL3049Meta.Severity,
				Message:  fmt.Sprintf("Label `%s` is missing.", requiredLabel),
				Line:     0,
			})
		}
	}

	return state
}
