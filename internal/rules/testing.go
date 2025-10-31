package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/internal/rule"
)

// LintDockerfile lints a Dockerfile string with the given rules and returns violations.
func LintDockerfile(dockerfile string, rules []rule.Rule) []rule.CheckFailure {
	p := parser.NewBuildkitParser()

	instructions, err := p.Parse([]byte(dockerfile))
	if err != nil {
		// Return empty if parse fails - some tests intentionally use invalid syntax
		return nil
	}

	processor := process.NewProcessor(rules)

	return processor.Run(instructions)
}

// AssertContainsViolation asserts that violations contains a failure with the given rule code.
func AssertContainsViolation(t *testing.T, violations []rule.CheckFailure, ruleCode string) {
	t.Helper()

	for _, v := range violations {
		if string(v.Code) == ruleCode {
			return
		}
	}

	t.Errorf("Expected violation %s not found. Got violations: %v", ruleCode, violations)
}

// AssertNoViolation asserts that violations does NOT contain the given rule code.
func AssertNoViolation(t *testing.T, violations []rule.CheckFailure, ruleCode string) {
	t.Helper()

	for _, v := range violations {
		if string(v.Code) == ruleCode {
			t.Errorf("Unexpected violation %s found: %+v", ruleCode, v)

			return
		}
	}
}
