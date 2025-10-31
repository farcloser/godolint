package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3061 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3061Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3061(t *testing.T) {
	allRules := []rule.Rule{ DL3061() }


	t.Run("don't warn: ARG then FROM then LABEL", func(t *testing.T) {
		dockerfile := `ARG A=B
FROM foo
LABEL foo=bar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3061")

	})

	t.Run("don't warn: FROM then ARG then RUN", func(t *testing.T) {
		dockerfile := `FROM foo
ARG A=B
RUN echo bla`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3061")

	})

	t.Run("don't warn: from before label", func(t *testing.T) {
		dockerfile := `FROM foo
LABEL foo=bar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3061")

	})

	t.Run("warn: label before from", func(t *testing.T) {
		dockerfile := `LABEL foo=bar
FROM foo`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3061")

	})

}
