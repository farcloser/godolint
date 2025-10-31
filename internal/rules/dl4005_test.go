package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL4005 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL4005Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL4005(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL4005(),
	}

	t.Run(
		"RUN ln",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN ln -sfv /bin/bash /bin/sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL4005")
		},
	)

	t.Run(
		"RUN ln with multiple acceptable commands",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN ln -s foo bar && unrelated && something_with /bin/sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL4005")
		},
	)

	t.Run(
		"RUN ln with unrelated symlinks",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN ln -sf /bin/true /sbin/initctl`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL4005")
		},
	)
}
