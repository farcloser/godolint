package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3058 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3058Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3058(t *testing.T) {
	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"maintainer": config.LabelTypeEmail,
		},
	}
	allRules := []rule.Rule{DL3058WithConfig(cfg)}

	t.Run("not ok with label not containing valid email", func(t *testing.T) {
		dockerfile := `LABEL maintainer="not-email"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3058")
	})

	t.Run("ok with label containing valid email", func(t *testing.T) {
		dockerfile := `LABEL maintainer="abcd@google.com"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3058")
	})

	t.Run("ok with other label not containing valid email", func(t *testing.T) {
		dockerfile := `LABEL other="doo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3058")
	})
}
