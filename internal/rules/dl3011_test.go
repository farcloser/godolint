package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3011 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3011Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3011(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3011(),
	}

	t.Run(
		"invalid port",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `EXPOSE 80000`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3011")
		},
	)

	t.Run(
		"invalid port in range",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `EXPOSE 40000-80000/tcp`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3011")
		},
	)

	t.Run(
		"valid port",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `EXPOSE 60000`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3011")
		},
	)

	t.Run(
		"valid port range",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `EXPOSE 40000-60000/tcp`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3011")
		},
	)

	t.Run(
		"valid port range variable",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `EXPOSE 40000-${FOOBAR}`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3011")
		},
	)

	t.Run(
		"valid port variable",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `EXPOSE ${FOOBAR}`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3011")
		},
	)
}
