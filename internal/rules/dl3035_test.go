package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3035 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3035Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3035(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3035(),
	}

	t.Run(
		"not ok: zypper dist-upgrade",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN zypper dist-upgrade`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3035")
		},
	)

	t.Run(
		"not ok: zypper dist-upgrade (2)",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN zypper dup`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3035")
		},
	)
}
