package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3015 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3015Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3015(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3015(),
	}

	t.Run(
		"apt-get no install recommends",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN apt-get install python`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3015")
		},
	)

	t.Run(
		"apt-get no install recommends",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN apt-get -y install python`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3015")
		},
	)

	t.Run(
		"apt-get no install recommends via option",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN apt-get -o APT::Install-Recommends=false install python`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3015")
		},
	)
}
