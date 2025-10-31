package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3021 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3021Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3021(t *testing.T) {
	allRules := []rule.Rule{ DL3021() }


	t.Run("no warn on 2 args", func(t *testing.T) {
		dockerfile := `COPY foo bar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3021")

	})

	t.Run("no warn on 3 args", func(t *testing.T) {
		dockerfile := `COPY foo bar baz/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3021")

	})

	t.Run("no warn on 3 args with quotes", func(t *testing.T) {
		dockerfile := `COPY foo bar "baz/"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3021")

	})

	t.Run("warn on 3 args", func(t *testing.T) {
		dockerfile := `COPY foo bar baz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3021")

	})

	t.Run("warn on 3 args with quotes", func(t *testing.T) {
		dockerfile := `COPY foo bar "baz"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3021")

	})

}
