package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3022 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3022Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3022(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3022(),
	}

	t.Run(
		"don't warn on correctly defined aliases",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM scratch as build
RUN foo
FROM node
COPY --from=build foo .
RUN baz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3022")
		},
	)

	t.Run(
		"don't warn on external images",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `COPY --from=haskell:latest bar .`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3022")
		},
	)

	t.Run(
		"don't warn on valid stage count as --from=<count>",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM scratch as build
RUN foo
FROM node
COPY --from=0 foo .`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3022")
		},
	)

	t.Run(
		"don't warn on valid stage count as --from=<count>",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM scratch
RUN foo
FROM node
COPY --from=0 foo .`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3022")
		},
	)

	t.Run(
		"warn on alias defined after",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM scratch
COPY --from=build foo .
FROM node as build
RUN baz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3022")
		},
	)

	t.Run(
		"warn on invalid stage count as --from=<count>",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `COPY --from=0 bar .`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3022")
		},
	)

	t.Run(
		"warn on missing alias",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `COPY --from=foo bar .`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3022")
		},
	)
}
