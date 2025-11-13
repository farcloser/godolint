package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3033 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3033Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3033(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3033(),
	}

	t.Run(
		"not ok without yum version pinning",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y tomcat && yum clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"not ok without yum version pinning - modules",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum module install -y tomcat && yum clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y tomcat-9.2 && yum clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning - modules",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum module install -y tomcat:9 && yum clean all`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning - modules (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning - package name contains `-`",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y rpm-sign-4.16.1.3`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning - package name contains `-` and `+`",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y gcc-c++-1.1.1`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)

	t.Run(
		"ok with yum version pinning - version contains epoch",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN yum install -y openssl-1:1.1.1k`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3033")
		},
	)
}
