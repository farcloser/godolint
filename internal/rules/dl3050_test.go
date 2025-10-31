package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3050 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3050Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3050(t *testing.T) {
	allRules := []rule.Rule{ DL3050() }


	t.Run("not ok with just other label", func(t *testing.T) {
		dockerfile := `LABEL other="bar"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3050")

	})

	t.Run("not ok with other label and required label", func(t *testing.T) {
		dockerfile := `LABEL required="foo" other="bar"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3050")

	})

	t.Run("ok with no label", func(t *testing.T) {
		dockerfile := ``
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3050")

	})

	t.Run("ok with required label", func(t *testing.T) {
		dockerfile := `LABEL required="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3050")

	})

}
