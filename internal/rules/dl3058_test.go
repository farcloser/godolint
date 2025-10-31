package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3058 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3058Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3058(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"maintainer": config.LabelTypeEmail,
		},
	}
	allRules := []rule.Rule{
		rules.DL3058WithConfig(cfg),
	}

	t.Run(
		"not ok with label not containing valid email",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL maintainer="not-email"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3058")
		},
	)

	t.Run(
		"ok with label containing valid email",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL maintainer="abcd@google.com"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3058")
		},
	)

	t.Run(
		"ok with other label not containing valid email",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `LABEL other="doo"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3058")
		},
	)
}
