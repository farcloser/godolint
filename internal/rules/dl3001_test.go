package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3001 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3001Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3001(t *testing.T) {
	allRules := []rule.Rule{DL3001()}

	t.Run("install ssh", func(t *testing.T) {
		dockerfile := `RUN apt-get install ssh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3001")
	})

	t.Run("invalid cmd", func(t *testing.T) {
		dockerfile := `RUN top`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3001")
	})
}
