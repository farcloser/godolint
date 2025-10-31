package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3006 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3006Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3006(t *testing.T) {
	allRules := []rule.Rule{ DL3006() }


	t.Run("no untagged", func(t *testing.T) {
		dockerfile := `FROM debian`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3006")

	})

	t.Run("no untagged with name", func(t *testing.T) {
		dockerfile := `FROM debian AS builder`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3006")

	})

	t.Run("scratch", func(t *testing.T) {
		dockerfile := `FROM scratch`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")

	})

	t.Run("untagged digest is not an error", func(t *testing.T) {
		dockerfile := `FROM ruby@sha256:f1dbca0f5dbc9`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")

	})

	t.Run("untagged digest is not an error", func(t *testing.T) {
		dockerfile := `FROM ruby:2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")

	})

	t.Run("using args is not an error", func(t *testing.T) {
		dockerfile := `FROM ${VALUE}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")

	})

}
