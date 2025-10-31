package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3022 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3022Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3022(t *testing.T) {
	allRules := []rule.Rule{ DL3022() }


	t.Run("don't warn on correctly defined aliases", func(t *testing.T) {
		dockerfile := `don't warn on correctly defined aliases
FROM scratch as build
RUN foo
FROM node
COPY --from=build foo .
RUN baz
DL3022`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3022")

	})

	t.Run("don't warn on external images", func(t *testing.T) {
		dockerfile := `COPY --from=haskell:latest bar .`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3022")

	})

	t.Run("don't warn on valid stage count as --from=<count>", func(t *testing.T) {
		dockerfile := `don't warn on valid stage count as --from=<count>
FROM scratch as build
RUN foo
FROM node
COPY --from=0 foo .
DL3022`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3022")

	})

	t.Run("don't warn on valid stage count as --from=<count>", func(t *testing.T) {
		dockerfile := `don't warn on valid stage count as --from=<count>
FROM scratch
RUN foo
FROM node
COPY --from=0 foo .
DL3022`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3022")

	})

	t.Run("warn on alias defined after", func(t *testing.T) {
		dockerfile := `warn on alias defined after
FROM scratch
COPY --from=build foo .
FROM node as build
RUN baz
DL3022`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3022")

	})

	t.Run("warn on invalid stage count as --from=<count>", func(t *testing.T) {
		dockerfile := `COPY --from=0 bar .`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3022")

	})

	t.Run("warn on missing alias", func(t *testing.T) {
		dockerfile := `COPY --from=foo bar .`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3022")

	})

}
