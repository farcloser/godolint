package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3061 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3061Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3061(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3061(),
	}

	t.Run(
		"don't warn: ARG then FROM then LABEL",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `ARG A=B
FROM foo
LABEL foo=bar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3061")
		},
	)

	t.Run(
		"don't warn: FROM then ARG then RUN",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM foo
ARG A=B
RUN echo bla`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3061")
		},
	)

	t.Run(
		"don't warn: from before label",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM foo
LABEL foo=bar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3061")
		},
	)

	t.Run(
		"warn: label before from",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL foo=bar
FROM foo`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3061")
		},
	)
}
