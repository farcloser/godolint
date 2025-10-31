package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3037 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3037Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3037(t *testing.T) {
	allRules := []rule.Rule{DL3037()}

	t.Run("not ok without zypper version pinning", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3037")
	})

	t.Run("ok with different variants of zypper version pinning", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat=9.0.39 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3037")
	})

	t.Run("ok with different variants of zypper version pinning (2)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat\>=9.0 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3037")
	})

	t.Run("ok with different variants of zypper version pinning (3)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat\>9.0 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3037")
	})

	t.Run("ok with different variants of zypper version pinning (4)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat\<=9.0 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3037")
	})

	t.Run("ok with different variants of zypper version pinning (5)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat\<9.0 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3037")
	})

	t.Run("ok with different variants of zypper version pinning (6)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y tomcat-9.0.39-1.rpm && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3037")
	})
}
