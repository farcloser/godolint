package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3027 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3027Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3027(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3027(),
	}

	t.Run(
		"apt",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM ubuntu
RUN apt install python`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3027")
		},
	)
}
