package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3029 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3029Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3029(t *testing.T) {
	allRules := []rule.Rule{ DL3029() }


	t.Run("allows platform $BUILDPLATFORM flag", func(t *testing.T) {
		dockerfile := `FROM --platform=$BUILDPLATFORM debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3029")

	})

	t.Run("allows platform $TARGETPLATFORM flag", func(t *testing.T) {
		dockerfile := `FROM --platform=$TARGETPLATFORM debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3029")

	})

	t.Run("allows platform ${BUILDPLATFORM:-} flag", func(t *testing.T) {
		dockerfile := `FROM --platform=${BUILDPLATFORM:-} debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3029")

	})

	t.Run("allows platform ${BUILDPLATFORM} flag", func(t *testing.T) {
		dockerfile := `FROM --platform=${BUILDPLATFORM} debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3029")

	})

	t.Run("explicit platform flag", func(t *testing.T) {
		dockerfile := `FROM --platform=linux debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3029")

	})

	t.Run("no platform flag", func(t *testing.T) {
		dockerfile := `FROM debian:jessie`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3029")

	})

}
