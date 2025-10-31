package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3047 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3047Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3047(t *testing.T) {
	allRules := []rule.Rule{DL3047()}

	t.Run(
		"does not warn when running with --append-output (append-output long option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget --append-output=/tmp/wget.log my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run(
		"does not warn when running with --no-verbose (no-verbose long option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget --no-verbose my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run(
		"does not warn when running with --output-file (output-file long option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget --output-file=/tmp/wget.log my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run("does not warn when running with --progress option", func(t *testing.T) {
		dockerfile := `FROM node as foo
RUN wget --progress=dot:giga my.xyz
DL3047`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3047")
	})

	t.Run(
		"does not warn when running with --quiet (quiet long option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget --quiet my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run(
		"does not warn when running with -a (append-output long option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget -a /tmp/wget.log my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run(
		"does not warn when running with -nv (no-verbose short option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget -nv my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run(
		"does not warn when running with -o (output-file long option) and without --progress option",
		func(t *testing.T) {
			dockerfile := `FROM node as foo
RUN wget -o /tmp/wget.log my.xyz
DL3047`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3047")
		},
	)

	t.Run("does not warn when running with -q (quiet short option) and without --progress option", func(t *testing.T) {
		dockerfile := `FROM node as foo
RUN wget -q my.xyz
DL3047`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3047")
	})

	t.Run("warns when using wget without --progress option", func(t *testing.T) {
		dockerfile := `FROM node as foo
RUN wget my.xyz
DL3047`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3047")
	})
}
