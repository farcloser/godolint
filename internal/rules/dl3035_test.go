package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3035 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3035Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3035(t *testing.T) {
	allRules := []rule.Rule{ DL3035() }


	t.Run("not ok: zypper dist-upgrade", func(t *testing.T) {
		dockerfile := `RUN zypper dist-upgrade`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3035")

	})

	t.Run("not ok: zypper dist-upgrade (2)", func(t *testing.T) {
		dockerfile := `RUN zypper dup`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3035")

	})

}
