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

// DL3054WithConfig creates the rule with custom configuration.
func DL3054WithConfig(cfg *config.Config) rule.Rule {
	return &DL3054Rule{
		cfg: cfg,
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
// This validates against common SPDX license identifiers and expressions.
func isValidSPDX(license string) bool {
	// Common SPDX license identifiers (not exhaustive, but covers common cases)
	commonLicenses := map[string]bool{
		"MIT": true, "Apache-2.0": true, "GPL-2.0": true, "GPL-3.0": true,
		"BSD-2-Clause": true, "BSD-3-Clause": true, "ISC": true,
		"LGPL-2.1": true, "LGPL-3.0": true, "MPL-2.0": true,
		"AGPL-3.0": true, "Unlicense": true, "CC0-1.0": true,
		"GPL-3.0-or-later": true, "GPL-2.0-or-later": true,
		"LGPL-3.0-or-later": true, "LGPL-2.1-or-later": true,
	}

	// Check if it's a simple license identifier
	if commonLicenses[license] {
		return true
	}

	// Check if it's a compound expression (e.g., "MIT AND Apache-2.0")
	// Pattern: valid-id (AND|OR|WITH valid-id)*
	exprPattern := regexp.MustCompile(`^[A-Za-z0-9.\-+()]+(?: (?:AND|OR|WITH) [A-Za-z0-9.\-+()]+)*$`)
	if !exprPattern.MatchString(license) {
		return false
	}

	// For compound expressions, validate each component
	parts := regexp.MustCompile(`\s+(?:AND|OR|WITH)\s+`).Split(license, -1)
	for _, part := range parts {
		// Check if part is a known license or follows SPDX naming convention
		if !commonLicenses[part] && !isSPDXPattern(part) {
			return false
		}
	}

	return true
}

// isSPDXPattern checks if a string follows common SPDX naming patterns.
func isSPDXPattern(s string) bool {
	// Common patterns: XXX-N.N, XXX-N.N-or-later, XXX-N.N-only
	patterns := []string{
		`^[A-Z][A-Za-z0-9]+-\d+\.\d+$`,                    // MIT-1.0
		`^[A-Z][A-Za-z0-9]+-\d+\.\d+-or-later$`,          // GPL-3.0-or-later
		`^[A-Z][A-Za-z0-9]+-\d+\.\d+-only$`,              // GPL-3.0-only
		`^[A-Z][A-Za-z0-9]+-\d+\.\d+-Clause$`,            // BSD-3-Clause
		`^CC-BY(-[A-Z]+)*-\d+\.\d+$`,                      // CC-BY-SA-4.0
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, s); matched {
			return true
		}
	}

	return false
}
