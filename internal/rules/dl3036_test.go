package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3036 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3036Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3036(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3036(),
	}

	t.Run(
		"not ok without zypper clean",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN zypper install -y mariadb=10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3036")
		},
	)

	t.Run(
		"ok when mount type cache is used",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=cache,target=/var/cache/zypp zypper install -y mariadb`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3036")
		},
	)

	t.Run(
		"ok when mount type tmpfs is used",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=tmpfs,target=/var/cache/zypp zypper install -y mariadb`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3036")
		},
	)

	t.Run(
		"ok with zypper clean",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN zypper install -y mariadb=10.4 && zypper clean`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3036")
		},
	)

	t.Run(
		"ok with zypper clean (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN zypper install -y mariadb=10.4 && zypper cc`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3036")
		},
	)
}
