package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3012 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3012Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3012(t *testing.T) {
	allRules := []rule.Rule{ DL3012() }


	t.Run("ok with no HEALTHCHECK instruction", func(t *testing.T) {
		dockerfile := `FROM scratch`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3012")

	})

	t.Run("ok with one HEALTHCHECK instruction", func(t *testing.T) {
		dockerfile := `FROM scratch
HEALTHCHECK CMD /bin/bla`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3012")

	})

	t.Run("ok with two HEALTHCHECK instructions in two stages", func(t *testing.T) {
		dockerfile := `FROM scratch
HEALTHCHECK CMD /bin/bla1
FROM scratch
HEALTHCHECK CMD /bin/bla2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3012")

	})

	t.Run("warn with two HEALTHCHECK instructions", func(t *testing.T) {
		dockerfile := `FROM scratch
HEALTHCHECK CMD /bin/bla1
HEALTHCHECK CMD /bin/bla2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3012")

	})

}
