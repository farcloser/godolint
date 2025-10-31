package rules

import (
	"fmt"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3051Rule checks that labels are not empty.
type DL3051Rule struct {
	cfg *config.Config
}

// DL3051 creates the rule for checking empty labels.
func DL3051() rule.Rule {
	return &DL3051Rule{
		cfg: config.Default(),
	}
}

// DL3051WithConfig creates the rule with custom configuration.
func DL3051WithConfig(cfg *config.Config) rule.Rule {
	return &DL3051Rule{
		cfg: cfg,
	}
}

func (r *DL3051Rule) Code() rule.RuleCode {
	return DL3051Meta.Code
}

func (r *DL3051Rule) Severity() rule.Severity {
	return DL3051Meta.Severity
}

func (r *DL3051Rule) Message() string {
	return DL3051Meta.Message
}

func (r *DL3051Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

func (r *DL3051Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Check each label in the schema
	for _, pair := range label.Pairs {
		// Only validate if this label is in the schema
		if _, exists := r.cfg.LabelSchema[pair.Key]; !exists {
			continue
		}

		// Check if value is empty
		if pair.Value == "" {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3051Meta.Code,
				Severity: DL3051Meta.Severity,
				Message:  fmt.Sprintf("label `%s` is empty.", pair.Key),
				Line:     line,
			})
		}
	}

	return state
}

func (r *DL3051Rule) Finalize(state rule.State) rule.State {
	return state
}
