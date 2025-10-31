package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/config"
	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3055 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3055Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3055(t *testing.T) {
	// Config: labelSchema = {"githash": GitHash}
	cfg := &config.Config{
		LabelSchema: map[string]config.LabelType{"githash": config.LabelTypeGitHash},
	}
	allRules := []rule.Rule{DL3055WithConfig(cfg)}


	t.Run("not ok with label not containing git hash", func(t *testing.T) {
		dockerfile := `LABEL githash="not-git-hash"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3055")

	})

	t.Run("ok with label containing long git hash", func(t *testing.T) {
		dockerfile := `LABEL githash="43c572f1272b6b3171dd1db9e41b7027128ce080"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3055")

	})

	t.Run("ok with label containing short git hash", func(t *testing.T) {
		dockerfile := `LABEL githash="2dbfae9"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3055")

	})

	t.Run("ok with other label not containing git hash", func(t *testing.T) {
		dockerfile := `LABEL other="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3055")

	})

}
