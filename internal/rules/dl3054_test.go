package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3054 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3054Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3054(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"spdxlabel": config.LabelTypeSPDX,
		},
	}
	allRules := []rule.Rule{
		rules.DL3054WithConfig(cfg),
	}

	t.Run(
		"not ok with label not containing SPDX identifier",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL spdxlabel="not-spdx"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3054")
		},
	)

	t.Run(
		"ok with label containing SPDX identifier",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL spdxlabel="BSD-3-Clause"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3054")
		},
	)

	t.Run(
		"ok with other label not containing SPDX identifier",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL other="fooo"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3054")
		},
	)
}
