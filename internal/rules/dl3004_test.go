package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3004 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3004Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3004(t *testing.T) {
	allRules := []rule.Rule{ DL3004() }


	t.Run("install sudo", func(t *testing.T) {
		dockerfile := `RUN apt-get install sudo`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3004")

	})

	t.Run("sudo", func(t *testing.T) {
		dockerfile := `RUN sudo apt-get update`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3004")

	})

	t.Run("sudo chained programs", func(t *testing.T) {
		dockerfile := `RUN apt-get update && sudo apt-get install`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3004")

	})

}
