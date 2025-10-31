package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3003 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3003Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3003(t *testing.T) {
	allRules := []rule.Rule{DL3003()}

	t.Run("not ok using cd", func(t *testing.T) {
		dockerfile := `RUN cd /opt`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3003")
	})

	t.Run("ok using WORKDIR", func(t *testing.T) {
		dockerfile := `WORKDIR /opt`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3003")
	})
}
