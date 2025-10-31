package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3060 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3060Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3060(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3060(),
	}

	t.Run(
		"not ok when cache mount is in wrong location",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=cache,target=/var/lib/foobar yarn install foobar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3060")
		},
	)

	t.Run(
		"not ok when tmpfs mount is in wrong location",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=tmpfs,target=/var/lib/foobar yarn install foobar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3060")
		},
	)

	t.Run(
		"not ok with no cache clean",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yarn install foo`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3060")
		},
	)

	t.Run(
		"ok when cache mount is used",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=cache,target=/root/.cache/yarn yarn install foobar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3060")
		},
	)

	t.Run(
		"ok when tmpfs mount is used",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=tmpfs,target=/root/.cache/yarn yarn install foobar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3060")
		},
	)

	t.Run(
		"ok with cache clean",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yarn install bar && yarn cache clean`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3060")
		},
	)

	t.Run(
		"ok with non-yarn commands",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN foo`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3060")
		},
	)
}
