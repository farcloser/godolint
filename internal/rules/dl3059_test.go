package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3059 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3059Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3059(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3059(),
	}

	t.Run(
		"not ok with two `RUN`s separated by a comment",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN /foo.sh
# a comment
RUN /bar.sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3059")
		},
	)

	t.Run(
		"not ok with two consecutive `RUN`s",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN /foo.sh
RUN /bar.sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3059")
		},
	)

	t.Run(
		"ok with no `RUN` at all",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian:10`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3059")
		},
	)

	t.Run(
		"ok with one `RUN`",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN /foo.sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3059")
		},
	)

	t.Run(
		"ok with one `RUN` after a comment",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `# a comment
RUN /foo.sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3059")
		},
	)

	t.Run(
		"ok with two not consecutive `RUN`",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN /foo.sh
WORKDIR /
RUN /bar.sh`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3059")
		},
	)
}
