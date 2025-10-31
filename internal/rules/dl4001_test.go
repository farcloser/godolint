package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL4001 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL4001Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL4001(t *testing.T) {
	allRules := []rule.Rule{ DL4001() }


	t.Run("does not warn when using both curl and wget in different stages", func(t *testing.T) {
		dockerfile := `does not warn when using both curl and wget in different stages
FROM node as foo
RUN wget my.xyz
FROM scratch
RUN curl localhost
DL4001`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4001")

	})

	t.Run("does not warn when using only wget", func(t *testing.T) {
		dockerfile := `does not warn when using only wget
FROM node as foo
RUN wget my.xyz
DL4001`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4001")

	})

	t.Run("does not warns when using both, on a single stage", func(t *testing.T) {
		dockerfile := `does not warns when using both, on a single stage
FROM node as foo
RUN wget my.xyz
RUN curl localhost
FROM scratch
RUN curl localhost
DL4001`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4001")

	})

	t.Run("warns when using both wget and curl", func(t *testing.T) {
		dockerfile := `warns when using both wget and curl
FROM node as foo
RUN wget my.xyz
RUN curl localhost
DL4001`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4001")

	})

	t.Run("warns when using both wget and curl in same instruction", func(t *testing.T) {
		dockerfile := `warns when using both wget and curl in same instruction
FROM node as foo
RUN wget my.xyz && curl localhost
DL4001`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4001")

	})

}
