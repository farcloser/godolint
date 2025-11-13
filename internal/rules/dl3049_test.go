package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3049 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3049Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3049(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"foo": config.LabelTypeRawText,
		},
	}
	allRules := []rule.Rule{
		rules.DL3049WithConfig(cfg),
	}

	t.Run(
		"not ok: single stage, no label",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM baseimage`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3049")
		},
	)

	t.Run(
		"not ok: single stage, wrong label",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM baseimage
LABEL bar="baz"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3049")
		},
	)

	t.Run(
		"ok: single stage, label present",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM baseimage
LABEL foo="bar"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3049")
		},
	)
}
