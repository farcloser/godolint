package rules

import (
	"fmt"
	"time"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3053Rule checks that RFC3339 timestamp labels are valid.
type DL3053Rule struct {
	cfg *config.Config
}

// DL3053 creates the rule for checking RFC3339 labels.
func DL3053() rule.Rule {
	return &DL3053Rule{
		cfg: config.Default(),
	}
}

// DL3053WithConfig creates the rule with custom configuration.
func DL3053WithConfig(cfg *config.Config) rule.Rule {
	return &DL3053Rule{
		cfg: cfg,
	}
}

func (r *DL3053Rule) Code() rule.RuleCode {
	return DL3053Meta.Code
}

func (r *DL3053Rule) Severity() rule.Severity {
	return DL3053Meta.Severity
}

func (r *DL3053Rule) Message() string {
	return DL3053Meta.Message
}

func (r *DL3053Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

func (r *DL3053Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Check each label that should be RFC3339
	for _, pair := range label.Pairs {
		labelType, exists := r.cfg.LabelSchema[pair.Key]
		if !exists || labelType != config.LabelTypeRFC3339 {
			continue
		}

		// Validate RFC3339 timestamp
		if !isValidRFC3339(pair.Value) {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3053Meta.Code,
				Severity: DL3053Meta.Severity,
				Message:  fmt.Sprintf("Label `%s` is not a valid time format - must conform to RFC3339.", pair.Key),
				Line:     line,
			})
		}
	}

	return state
}

func (r *DL3053Rule) Finalize(state rule.State) rule.State {
	return state
}

// isValidRFC3339 checks if a string is a valid RFC3339 timestamp.
func isValidRFC3339(timestamp string) bool {
	_, err := time.Parse(time.RFC3339, timestamp)
	return err == nil
}
