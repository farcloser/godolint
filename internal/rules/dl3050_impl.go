package rules

import (
	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3050Rule checks for superfluous labels (labels not in schema).
type DL3050Rule struct {
	cfg *config.Config
}

// DL3050 creates a rule that checks for superfluous labels.
// TODO: Wire up config from main application.
func DL3050() rule.Rule {
	return &DL3050Rule{
		cfg: config.Default(), // Use default config for now
	}
}

// DL3050WithConfig creates the rule with custom configuration.
func DL3050WithConfig(cfg *config.Config) rule.Rule {
	return &DL3050Rule{
		cfg: cfg,
	}
}

// Code returns the rule code.
func (*DL3050Rule) Code() rule.RuleCode {
	return DL3050Meta.Code
}

// Severity returns the rule severity.
func (*DL3050Rule) Severity() rule.Severity {
	return DL3050Meta.Severity
}

// Message returns the rule message.
func (*DL3050Rule) Message() string {
	return DL3050Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3050Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

// Check validates that LABEL instructions only use labels defined in schema.
func (r *DL3050Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Only check if strict labels is enabled
	if !r.cfg.StrictLabels {
		return state
	}

	// Check each label against schema
	for _, pair := range label.Pairs {
		if _, exists := r.cfg.LabelSchema[pair.Key]; !exists {
			// Label not in schema
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3050Meta.Code,
				Severity: DL3050Meta.Severity,
				Message:  DL3050Meta.Message,
				Line:     line,
			})
		}
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*DL3050Rule) Finalize(state rule.State) rule.State {
	return state
}
