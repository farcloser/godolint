package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3034 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3034Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3034(t *testing.T) {
	allRules := []rule.Rule{ DL3034() }


	t.Run("not ok without non-interactive switch", func(t *testing.T) {
		dockerfile := `RUN zypper install httpd=2.4.24 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3034")

	})

	t.Run("ok with non-interactive switch present", func(t *testing.T) {
		dockerfile := `RUN zypper install -n httpd=2.4.24 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3034")

	})

	t.Run("ok with non-interactive switch present (2)", func(t *testing.T) {
		dockerfile := `RUN zypper install --non-interactive httpd=2.4.24 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3034")

	})

	t.Run("ok with non-interactive switch present (3)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y httpd=2.4.24 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3034")

	})

	t.Run("ok with non-interactive switch present (4)", func(t *testing.T) {
		dockerfile := `RUN zypper install --no-confirm httpd=2.4.24 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3034")

	})

}
