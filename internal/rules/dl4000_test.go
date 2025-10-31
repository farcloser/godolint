package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL4000 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL4000Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL4000(t *testing.T) {
	allRules := []rule.Rule{DL4000()}

	t.Run("has deprecated maintainer", func(t *testing.T) {
		dockerfile := `FROM busybox
MAINTAINER hudu@mail.com`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4000")
	})

	t.Run("has maintainer", func(t *testing.T) {
		dockerfile := `FROM debian
MAINTAINER Lukas`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4000")
	})

	t.Run("has maintainer first", func(t *testing.T) {
		dockerfile := `MAINTAINER Lukas
FROM DEBIAN`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4000")
	})

	t.Run("has no maintainer", func(t *testing.T) {
		dockerfile := `FROM debian`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4000")
	})
}
