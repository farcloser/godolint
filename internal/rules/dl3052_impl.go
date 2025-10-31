package rules

import (
	"fmt"
	"net/url"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3052Rule checks that URL labels are valid.
type DL3052Rule struct {
	cfg *config.Config
}

// DL3052 creates the rule for checking URL labels.
func DL3052() rule.Rule {
	return &DL3052Rule{
		cfg: config.Default(),
	}
}

func (r *DL3052Rule) Code() rule.RuleCode {
	return DL3052Meta.Code
}

func (r *DL3052Rule) Severity() rule.Severity {
	return DL3052Meta.Severity
}

func (r *DL3052Rule) Message() string {
	return DL3052Meta.Message
}

func (r *DL3052Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

func (r *DL3052Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Check each label that should be a URL
	for _, pair := range label.Pairs {
		labelType, exists := r.cfg.LabelSchema[pair.Key]
		if !exists || labelType != config.LabelTypeURL {
			continue
		}

		// Validate URL
		if !isValidURL(pair.Value) {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3052Meta.Code,
				Severity: DL3052Meta.Severity,
				Message:  fmt.Sprintf("Label `%s` is not a valid URL.", pair.Key),
				Line:     line,
			})
		}
	}

	return state
}

func (r *DL3052Rule) Finalize(state rule.State) rule.State {
	return state
}

// isValidURL checks if a string is a valid URL.
func isValidURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	// Must have a scheme (http, https, etc.)
	return u.Scheme != ""
}
