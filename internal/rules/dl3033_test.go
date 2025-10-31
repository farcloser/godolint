package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3033 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3033Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3033(t *testing.T) {
	allRules := []rule.Rule{ DL3033() }


	t.Run("not ok without yum version pinning", func(t *testing.T) {
		dockerfile := `RUN yum install -y tomcat && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3033")

	})

	t.Run("not ok without yum version pinning - modules", func(t *testing.T) {
		dockerfile := `RUN yum module install -y tomcat && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning", func(t *testing.T) {
		dockerfile := `RUN yum install -y tomcat-9.2 && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning (2)", func(t *testing.T) {
		dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning - modules", func(t *testing.T) {
		dockerfile := `RUN yum module install -y tomcat:9 && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning - modules (2)", func(t *testing.T) {
		dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning - package name contains `-`", func(t *testing.T) {
		dockerfile := `RUN yum install -y rpm-sign-4.16.1.3`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning - package name contains `-` and `+`", func(t *testing.T) {
		dockerfile := `RUN yum install -y gcc-c++-1.1.1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

	t.Run("ok with yum version pinning - version contains epoch", func(t *testing.T) {
		dockerfile := `RUN yum install -y openssl-1:1.1.1k`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3033")

	})

}
