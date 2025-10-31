package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL4005 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL4005Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL4005(t *testing.T) {
	allRules := []rule.Rule{ DL4005() }


	t.Run("RUN ln", func(t *testing.T) {
		dockerfile := `RUN ln -sfv /bin/bash /bin/sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4005")

	})

	t.Run("RUN ln with multiple acceptable commands", func(t *testing.T) {
		dockerfile := `RUN ln -s foo bar && unrelated && something_with /bin/sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4005")

	})

	t.Run("RUN ln with unrelated symlinks", func(t *testing.T) {
		dockerfile := `RUN ln -sf /bin/true /sbin/initctl`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4005")

	})

}
