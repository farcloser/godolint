package rules

import (
	"fmt"
	"net/mail"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3058Rule checks that email labels are valid.
type DL3058Rule struct {
	cfg *config.Config
}

// DL3058 creates the rule for checking email labels.
func DL3058() rule.Rule {
	return &DL3058Rule{
		cfg: config.Default(),
	}
}

// DL3058WithConfig creates the rule with custom configuration.
func DL3058WithConfig(cfg *config.Config) rule.Rule {
	return &DL3058Rule{
		cfg: cfg,
	}
}

// Code returns the rule code.
func (*DL3058Rule) Code() rule.RuleCode {
	return DL3058Meta.Code
}

// Severity returns the rule severity.
func (*DL3058Rule) Severity() rule.Severity {
	return DL3058Meta.Severity
}

// Message returns the rule message.
func (*DL3058Rule) Message() string {
	return DL3058Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3058Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

// Check validates that email labels are valid email addresses.
func (r *DL3058Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Check each label that should be an email
	for _, pair := range label.Pairs {
		labelType, exists := r.cfg.LabelSchema[pair.Key]
		if !exists || labelType != config.LabelTypeEmail {
			continue
		}

		// Validate email
		if !isValidEmail(pair.Value) {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3058Meta.Code,
				Severity: DL3058Meta.Severity,
				Message:  fmt.Sprintf("Label `%s` is not a valid email format - must conform to RFC5322.", pair.Key),
				Line:     line,
			})
		}
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*DL3058Rule) Finalize(state rule.State) rule.State {
	return state
}

// isValidEmail checks if a string is a valid email address (RFC5322).
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}
