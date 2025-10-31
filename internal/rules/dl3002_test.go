package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3002 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3002Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3002(t *testing.T) {
	allRules := []rule.Rule{DL3002()}

	t.Run("can switch back to non root", func(t *testing.T) {
		dockerfile := `FROM scratch
USER root
RUN something
USER foo
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3002")
	})

	t.Run("does not warn when switching in multiple stages", func(t *testing.T) {
		dockerfile := `FROM debian as base
USER root
RUN something
USER foo
FROM scratch
RUN something else
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3002")
	})

	t.Run("last user should not be root", func(t *testing.T) {
		dockerfile := `FROM scratch
USER root
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3002")
	})

	t.Run("no UID:GID", func(t *testing.T) {
		dockerfile := `FROM scratch
USER 0:0
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3002")
	})

	t.Run("no root", func(t *testing.T) {
		dockerfile := `FROM scratch
USER foo
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3002")
	})

	t.Run("no root UID", func(t *testing.T) {
		dockerfile := `FROM scratch
USER 0
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3002")
	})

	t.Run("no root:root", func(t *testing.T) {
		dockerfile := `FROM scratch
USER root:root
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3002")
	})

	t.Run("warns on multiple stages", func(t *testing.T) {
		dockerfile := `FROM debian as base
USER root
RUN something
FROM scratch
USER foo
RUN something else
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3002")
	})

	t.Run("warns on transitive root user", func(t *testing.T) {
		dockerfile := `FROM debian as base
USER root
RUN something
FROM base
RUN something else
DL3002`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3002")
	})
}
