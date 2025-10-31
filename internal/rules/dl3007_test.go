package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3007 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3007Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3007(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3007(),
	}

	t.Run(
		"explicit latest",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian:latest`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3007")
		},
	)

	t.Run(
		"explicit latest with name",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian:latest AS builder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3007")
		},
	)

	t.Run(
		"explicit tagged",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian:jessie`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3007")
		},
	)

	t.Run(
		"explicit tagged with name",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian:jessie AS builder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3007")
		},
	)
}
