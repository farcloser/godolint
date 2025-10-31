package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3015 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3015Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3015(t *testing.T) {
	allRules := []rule.Rule{ DL3015() }


	t.Run("apt-get no install recommends", func(t *testing.T) {
		dockerfile := `RUN apt-get install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3015")

	})

	t.Run("apt-get no install recommends", func(t *testing.T) {
		dockerfile := `RUN apt-get -y install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3015")

	})

	t.Run("apt-get no install recommends via option", func(t *testing.T) {
		dockerfile := `RUN apt-get -o APT::Install-Recommends=false install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3015")

	})

}
