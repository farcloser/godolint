package rules_test

import (
	"testing"

	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/rules"
	"github.com/forkcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3066 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3066Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3066(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3066(),
	}

	t.Run(
		"not ok: non-numeric UID",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `USER foobar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3066")
		},
	)

	t.Run(
		"not ok: non-numeric UID and GID",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `USER foobar:barfoo`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3066")
		},
	)

	t.Run(
		"ok: UID is non-numeric, but is a defined ARG",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `ARG APP_UID
FROM foobar:barfoo
USER ${APP_UID}`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3066")
		},
	)

	t.Run(
		"ok: numeric UID",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `USER 12345`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3066")
		},
	)

	t.Run(
		"ok: numeric UID and GID",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `USER 1234:5678`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3066")
		},
	)
}
