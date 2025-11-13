package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3032 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3032Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3032(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3032(),
	}

	t.Run(
		"not ok with no clean all",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y mariadb-10.4`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3032")
		},
	)

	t.Run(
		"ok with rm -rf /var/cache/yum/*",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y mariadb-10.4 && rm -rf /var/cache/yum/*`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3032")
		},
	)

	t.Run(
		"ok with yum clean all ",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y mariadb-10.4 && yum clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3032")
		},
	)

	t.Run(
		"ok with yum clean all  (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3032")
		},
	)
}
