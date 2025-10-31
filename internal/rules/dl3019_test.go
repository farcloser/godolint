package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3019 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3019Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3019(t *testing.T) {
	allRules := []rule.Rule{DL3019()}

	t.Run("don't warn: apk add with --no-cache", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache flex=2.6.4-r1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3019")
	})

	t.Run("don't warn: apk add with BuildKit cache mount", func(t *testing.T) {
		dockerfile := `RUN --mount=type=cache,target=/var/cache/apk apk add -U curl=7.77.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3019")
	})

	t.Run("don't warn: apk add with BuildKit cache mount in wrong dir and --no-cache", func(t *testing.T) {
		dockerfile := `RUN --mount=type=cache,target=/var/cache/foo apk add --no-cache -U curl=7.77.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3019")
	})

	t.Run("don't warn: apk add with BuildKit tmpfs mount", func(t *testing.T) {
		dockerfile := `RUN --mount=type=tmpfs,target=/var/cache/apk apk add -U curl=7.77.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3019")
	})

	t.Run("warn: apk add with BuildKit cache mount to wrong dir", func(t *testing.T) {
		dockerfile := `RUN --mount=type=cache,target=/var/cache/foo apk add -U curl=7.77.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3019")
	})

	t.Run("warn: apk add with BuildKit tmpfs mount to wrong dir", func(t *testing.T) {
		dockerfile := `RUN --mount=type=tmpfs,target=/var/cache/foo apk add -U curl=7.77.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3019")
	})

	t.Run("warn: apk add without --no-cache", func(t *testing.T) {
		dockerfile := `RUN apk add flex=2.6.4-r1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3019")
	})
}
