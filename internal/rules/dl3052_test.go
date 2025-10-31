package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3052 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3052Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3052(t *testing.T) {
	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"urllabel": config.LabelTypeURL,
		},
	}
	allRules := []rule.Rule{DL3052WithConfig(cfg)}

	t.Run("not ok with label not containing URL", func(t *testing.T) {
		dockerfile := `LABEL urllabel="not-url"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3052")
	})

	t.Run("ok with label containing URL", func(t *testing.T) {
		dockerfile := `LABEL urllabel="http://example.com"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3052")
	})

	t.Run("ok with other label not containing URL", func(t *testing.T) {
		dockerfile := `LABEL other="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3052")
	})
}
