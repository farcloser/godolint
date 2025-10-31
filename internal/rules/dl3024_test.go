package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3024 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3024Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3024(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3024(),
	}

	t.Run(
		"don't warn on unique aliases",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM scratch as build
RUN foo
FROM node as run
RUN baz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3024")
		},
	)

	t.Run(
		"warn on duplicate aliases",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM node as foo
RUN something
FROM scratch as foo
RUN something`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3024")
		},
	)
}
