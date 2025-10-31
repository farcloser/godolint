package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3040 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3040Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3040(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3040(),
	}

	t.Run(
		"no ok without dnf clean all",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN dnf install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"no ok without dnf clean all (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN microdnf install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with cache mount at /var/cache/yum",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=cache,target=/var/cache/libdnf5 dnf install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with cache mount at /var/cache/yum (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=cache,target=/var/cache/libdnf5 microdnf install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with dnf clean all",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN dnf install -y mariadb-10.4 && dnf clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with dnf clean all (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN microdnf install -y mariadb-10.4 && microdnf clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with dnf clean all (3)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN notdnf install mariadb`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with rm /var/cache/yum",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN dnf install -y mariadb-10.4 && rm -rf /var/cache/libdnf5`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with rm /var/cache/yum (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN microdnf install -y mariadb-10.4 && rm -rf /var/cache/libdnf5`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with tmpfs mount at /var/cache/yum",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=tmpfs,target=/var/cache/libdnf5 dnf install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)

	t.Run(
		"ok with tmpfs mount at /var/cache/yum (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN --mount=type=tmpfs,target=/var/cache/libdnf5 microdnf install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3040")
		},
	)
}
