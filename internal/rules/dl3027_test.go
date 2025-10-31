package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3027 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3027Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3027(t *testing.T) {
	allRules := []rule.Rule{ DL3027() }


	t.Run("apt", func(t *testing.T) {
		dockerfile := `apt
FROM ubuntu
RUN apt install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3027")

	})

}
