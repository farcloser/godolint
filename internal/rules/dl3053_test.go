package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3053 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3053Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3053(t *testing.T) {
	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"datelabel": config.LabelTypeRFC3339,
		},
	}
	allRules := []rule.Rule{DL3053WithConfig(cfg)}

	t.Run("not ok with label not containing RFC3339 date", func(t *testing.T) {
		dockerfile := `LABEL datelabel="not-date"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3053")
	})

	t.Run("ok with label containing RFC3339 date", func(t *testing.T) {
		dockerfile := `LABEL datelabel="2021-03-10T10:26:33.564595127+01:00"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3053")
	})

	t.Run("ok with other label not containing RFC3339 date", func(t *testing.T) {
		dockerfile := `LABEL other="doo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3053")
	})
}
