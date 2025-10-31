package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3041 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3041Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3041(t *testing.T) {
	allRules := []rule.Rule{DL3041()}

	t.Run("not ok without dnf version pinning", func(t *testing.T) {
		dockerfile := `RUN dnf install -y tomcat && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3041")
	})

	t.Run("not ok without dnf version pinning (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y tomcat && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3041")
	})

	t.Run("not ok without dnf version pinning - modules", func(t *testing.T) {
		dockerfile := `RUN dnf module install -y tomcat && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3041")
	})

	t.Run("not ok without dnf version pinning - modules (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf module install -y tomcat && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3041")
	})

	t.Run("not ok without dnf version pinning - package name with `-`", func(t *testing.T) {
		dockerfile := `RUN dnf install -y rpm-sign && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3041")
	})

	t.Run("not ok without dnf version pinning - package name with `-` (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y rpm-sign && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning", func(t *testing.T) {
		dockerfile := `RUN dnf install -y tomcat-9.0.1 && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y tomcat-9.0.1 && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - modules", func(t *testing.T) {
		dockerfile := `RUN dnf module install -y tomcat:9 && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - modules (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf module install -y tomcat:9 && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - modules (3)", func(t *testing.T) {
		dockerfile := `RUN notdnf module install tomcat`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - package name with `-`", func(t *testing.T) {
		dockerfile := `RUN dnf install -y rpm-sign-4.16.1.3 && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - package name with `-` (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y rpm-sign-4.16.1.3 && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - package name with `-` and `+`", func(t *testing.T) {
		dockerfile := `RUN dnf install -y gcc-c++-1.1.1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - package name with `-` and `+` (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y gcc-c++-1.1.1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - package version with epoch", func(t *testing.T) {
		dockerfile := `RUN dnf install -y openssl-1:1.1.1k`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with dnf version pinning - package version with epoch (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y openssl-1:1.1.1k`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok with version pinning if command is not `dnf` or `microdnf`", func(t *testing.T) {
		dockerfile := `RUN notdnf install openssl-1:1.1.1k`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})

	t.Run("ok without version pinning if command is not `dnf` or `microdnf`", func(t *testing.T) {
		dockerfile := `RUN notdnf install tomcat`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3041")
	})
}
