package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3011 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3011Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3011(t *testing.T) {
	allRules := []rule.Rule{ DL3011() }


	t.Run("invalid port", func(t *testing.T) {
		dockerfile := `EXPOSE 80000`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3011")

	})

	t.Run("invalid port in range", func(t *testing.T) {
		dockerfile := `EXPOSE 40000-80000/tcp`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3011")

	})

	t.Run("valid port", func(t *testing.T) {
		dockerfile := `EXPOSE 60000`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3011")

	})

	t.Run("valid port range", func(t *testing.T) {
		dockerfile := `EXPOSE 40000-60000/tcp`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3011")

	})

	t.Run("valid port range variable", func(t *testing.T) {
		dockerfile := `EXPOSE 40000-${FOOBAR}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3011")

	})

	t.Run("valid port variable", func(t *testing.T) {
		dockerfile := `EXPOSE ${FOOBAR}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3011")

	})

}
