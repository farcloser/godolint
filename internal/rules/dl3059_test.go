package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3059 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3059Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3059(t *testing.T) {
	allRules := []rule.Rule{DL3059()}

	t.Run("not ok with two `RUN`s separated by a comment", func(t *testing.T) {
		dockerfile := `RUN /foo.sh
# a comment
RUN /bar.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3059")
	})

	t.Run("not ok with two consecutive `RUN`s", func(t *testing.T) {
		dockerfile := `RUN /foo.sh
RUN /bar.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3059")
	})

	t.Run("ok with no `RUN` at all", func(t *testing.T) {
		dockerfile := `FROM debian:10`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3059")
	})

	t.Run("ok with one `RUN`", func(t *testing.T) {
		dockerfile := `RUN /foo.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3059")
	})

	t.Run("ok with one `RUN` after a comment", func(t *testing.T) {
		dockerfile := `# a comment
RUN /foo.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3059")
	})

	t.Run("ok with two not consecutive `RUN`", func(t *testing.T) {
		dockerfile := `RUN /foo.sh
WORKDIR /
RUN /bar.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3059")
	})
}
