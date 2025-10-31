package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3046 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3046Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3046(t *testing.T) {
	allRules := []rule.Rule{DL3046()}

	t.Run("ok with `useradd` alone", func(t *testing.T) {
		dockerfile := `RUN useradd luser`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3046")
	})

	t.Run("ok with `useradd` and just flag `-l`", func(t *testing.T) {
		dockerfile := `RUN useradd -l luser`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3046")
	})

	t.Run("ok with `useradd` long uid and flag `-l`", func(t *testing.T) {
		dockerfile := `RUN useradd -l -u 123456 luser`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3046")
	})

	t.Run("ok with `useradd` short uid", func(t *testing.T) {
		dockerfile := `RUN useradd -u 12345 luser`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3046")
	})

	t.Run("warn when `useradd` and long uid without flag `-l`", func(t *testing.T) {
		dockerfile := `RUN useradd -u 123456 luser`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3046")
	})
}
