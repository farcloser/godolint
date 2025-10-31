package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3038 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3038Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3038(t *testing.T) {
	allRules := []rule.Rule{DL3038()}

	t.Run("not ok without dnf non-interactive flag", func(t *testing.T) {
		dockerfile := `RUN dnf install httpd-2.4.24 && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3038")
	})

	t.Run("not ok without dnf non-interactive flag (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install httpd-2.4.24 && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3038")
	})

	t.Run("ok with dnf non-interactive flag", func(t *testing.T) {
		dockerfile := `RUN dnf install -y httpd-2.4.24 && dnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3038")
	})

	t.Run("ok with dnf non-interactive flag (2)", func(t *testing.T) {
		dockerfile := `RUN microdnf install -y httpd-2.4.24 && microdnf clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3038")
	})

	t.Run("ok with dnf non-interactive flag (3)", func(t *testing.T) {
		dockerfile := `RUN notdnf install httpd`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3038")
	})
}
