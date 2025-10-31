package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3051 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3051Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3051(t *testing.T) {
	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{
			"emptylabel": config.LabelTypeRawText,
		},
	}
	allRules := []rule.Rule{DL3051WithConfig(cfg)}

	t.Run("not ok with label empty", func(t *testing.T) {
		dockerfile := `LABEL emptylabel=""`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3051")
	})

	t.Run("ok with label not empty", func(t *testing.T) {
		dockerfile := `LABEL emptylabel="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3051")
	})

	t.Run("ok with other label empty", func(t *testing.T) {
		dockerfile := `LABEL other=""`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3051")
	})
}
