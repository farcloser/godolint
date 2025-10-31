package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL4004 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL4004Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL4004(t *testing.T) {
	allRules := []rule.Rule{DL4004()}

	t.Run("many entrypoints", func(t *testing.T) {
		dockerfile := `FROM debian
ENTRYPOINT bash
RUN foo
ENTRYPOINT another
DL4004`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4004")
	})

	t.Run("many entrypoints, different stages", func(t *testing.T) {
		dockerfile := `FROM debian as distro1
ENTRYPOINT bash
RUN foo
ENTRYPOINT another
FROM debian as distro2
ENTRYPOINT another
DL4004`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4004")
	})

	t.Run("no cmd", func(t *testing.T) {
		dockerfile := `FROM busybox`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4004")
	})

	t.Run("no entry", func(t *testing.T) {
		dockerfile := `FROM busybox`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4004")
	})

	t.Run("single entry", func(t *testing.T) {
		dockerfile := `ENTRYPOINT /bin/true`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4004")
	})

	t.Run("single entrypoint, different stages", func(t *testing.T) {
		dockerfile := `FROM debian as distro1
ENTRYPOINT bash
RUN foo
FROM debian as distro2
ENTRYPOINT another
DL4004`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4004")
	})
}
