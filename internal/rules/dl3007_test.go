package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3007 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3007Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3007(t *testing.T) {
	allRules := []rule.Rule{DL3007()}

	t.Run("explicit latest", func(t *testing.T) {
		dockerfile := `FROM debian:latest`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3007")
	})

	t.Run("explicit latest with name", func(t *testing.T) {
		dockerfile := `FROM debian:latest AS builder`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3007")
	})

	t.Run("explicit tagged", func(t *testing.T) {
		dockerfile := `FROM debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3007")
	})

	t.Run("explicit tagged with name", func(t *testing.T) {
		dockerfile := `FROM debian:jessie AS builder`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3007")
	})
}
