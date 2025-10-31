package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3023 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3023Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3023(t *testing.T) {
	allRules := []rule.Rule{DL3023()}

	t.Run("don't warn on copying from other sources", func(t *testing.T) {
		dockerfile := `FROM scratch as build
RUN foo
FROM node as run
COPY --from=build foo .
RUN baz
DL3023`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3023")
	})

	t.Run("warn on copying from your the same FROM", func(t *testing.T) {
		dockerfile := `FROM node as foo
COPY --from=foo bar .
DL3023`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3023")
	})
}
