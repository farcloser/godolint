package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3054 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3054Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3054(t *testing.T) {
	allRules := []rule.Rule{ DL3054() }


	t.Run("not ok with label not containing SPDX identifier", func(t *testing.T) {
		dockerfile := `LABEL spdxlabel="not-spdx"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3054")

	})

	t.Run("ok with label containing SPDX identifier", func(t *testing.T) {
		dockerfile := `LABEL spdxlabel="BSD-3-Clause"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3054")

	})

	t.Run("ok with other label not containing SPDX identifier", func(t *testing.T) {
		dockerfile := `LABEL other="fooo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3054")

	})

}
