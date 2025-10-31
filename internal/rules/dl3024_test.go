package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3024 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3024Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3024(t *testing.T) {
	allRules := []rule.Rule{DL3024()}

	t.Run("don't warn on unique aliases", func(t *testing.T) {
		dockerfile := `FROM scratch as build
RUN foo
FROM node as run
RUN baz
DL3024`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3024")
	})

	t.Run("warn on duplicate aliases", func(t *testing.T) {
		dockerfile := `FROM node as foo
RUN something
FROM scratch as foo
RUN something
DL3024`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3024")
	})
}
