package rules

import (
	"fmt"
	"regexp"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3055Rule checks that git hash labels are valid.
type DL3055Rule struct {
	cfg *config.Config
}

// DL3055 creates the rule for checking git hash labels.
func DL3055() rule.Rule {
	return &DL3055Rule{
		cfg: config.Default(),
	}
}

// DL3055WithConfig creates the rule with custom configuration.
func DL3055WithConfig(cfg *config.Config) rule.Rule {
	return &DL3055Rule{
		cfg: cfg,
	}
}

func (r *DL3055Rule) Code() rule.RuleCode {
	return DL3055Meta.Code
}

func (r *DL3055Rule) Severity() rule.Severity {
	return DL3055Meta.Severity
}

func (r *DL3055Rule) Message() string {
	return DL3055Meta.Message
}

func (r *DL3055Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

func (r *DL3055Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Check each label that should be a git hash
	for _, pair := range label.Pairs {
		labelType, exists := r.cfg.LabelSchema[pair.Key]
		if !exists || labelType != config.LabelTypeGitHash {
			continue
		}

		// Validate git hash: must be 7 or 40 hex characters
		if !isValidGitHash(pair.Value) {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3055Meta.Code,
				Severity: DL3055Meta.Severity,
				Message:  fmt.Sprintf("Label `%s` is not a valid git hash.", pair.Key),
				Line:     line,
			})
		}
	}

	return state
}

func (r *DL3055Rule) Finalize(state rule.State) rule.State {
	return state
}

// isValidGitHash checks if a string is a valid git hash.
// Valid git hash: 7 or 40 hexadecimal characters.
func isValidGitHash(hash string) bool {
	if len(hash) != 7 && len(hash) != 40 {
		return false
	}

	// Check if all characters are hexadecimal
	hexPattern := regexp.MustCompile(`^[0-9a-fA-F]+$`)

	return hexPattern.MatchString(hash)
}
