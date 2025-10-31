package rules

import (
	"fmt"
	"regexp"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3054Rule checks that SPDX license labels are valid.
type DL3054Rule struct {
	cfg *config.Config
}

// DL3054 creates the rule for checking SPDX labels.
func DL3054() rule.Rule {
	return &DL3054Rule{
		cfg: config.Default(),
	}
}

func (r *DL3054Rule) Code() rule.RuleCode {
	return DL3054Meta.Code
}

func (r *DL3054Rule) Severity() rule.Severity {
	return DL3054Meta.Severity
}

func (r *DL3054Rule) Message() string {
	return DL3054Meta.Message
}

func (r *DL3054Rule) InitialState() rule.State {
	return rule.EmptyState(nil)
}

func (r *DL3054Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	label, ok := instruction.(*syntax.Label)
	if !ok {
		return state
	}

	// Check each label that should be an SPDX identifier
	for _, pair := range label.Pairs {
		labelType, exists := r.cfg.LabelSchema[pair.Key]
		if !exists || labelType != config.LabelTypeSPDX {
			continue
		}

		// Validate SPDX identifier
		if !isValidSPDX(pair.Value) {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL3054Meta.Code,
				Severity: DL3054Meta.Severity,
				Message:  fmt.Sprintf("Label `%s` is not a valid SPDX identifier.", pair.Key),
				Line:     line,
			})
		}
	}

	return state
}

func (r *DL3054Rule) Finalize(state rule.State) rule.State {
	return state
}

// isValidSPDX checks if a string is a valid SPDX license expression.
// This is a simplified check that validates basic SPDX patterns.
// For full SPDX validation, consider using github.com/spdx/tools-golang
func isValidSPDX(license string) bool {
	// Basic SPDX pattern: alphanumeric with dots, hyphens, plus signs
	// Examples: MIT, Apache-2.0, GPL-3.0-or-later, MIT AND Apache-2.0
	// Pattern allows: letters, numbers, dots, hyphens, plus, AND, OR, WITH
	pattern := regexp.MustCompile(`^[A-Za-z0-9.\-+()]+(?: (?:AND|OR|WITH) [A-Za-z0-9.\-+()]+)*$`)
	return pattern.MatchString(license)
}
