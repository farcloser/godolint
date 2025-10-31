package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL4003 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL4003Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL4003(t *testing.T) {
	allRules := []rule.Rule{ DL4003() }


	t.Run("many cmds", func(t *testing.T) {
		dockerfile := `many cmds
FROM debian
CMD bash
RUN foo
CMD another
DL4003`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4003")

	})

	t.Run("many cmds, different stages", func(t *testing.T) {
		dockerfile := `many cmds, different stages
FROM debian as distro1
CMD bash
RUN foo
CMD another
FROM debian as distro2
CMD another
DL4003`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL4003")

	})

	t.Run("single cmd", func(t *testing.T) {
		dockerfile := `CMD /bin/true`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4003")

	})

	t.Run("single cmds, different stages", func(t *testing.T) {
		dockerfile := `single cmds, different stages
FROM debian as distro1
CMD bash
RUN foo
FROM debian as distro2
CMD another
DL4003`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL4003")

	})

}
