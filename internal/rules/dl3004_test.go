package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3004 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3004Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3004(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3004(),
	}

	t.Run(
		"install sudo",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN apt-get install sudo`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3004")
		},
	)

	t.Run(
		"sudo",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN sudo apt-get update`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3004")
		},
	)

	t.Run(
		"sudo chained programs",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN apt-get update && sudo apt-get install`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3004")
		},
	)
}
