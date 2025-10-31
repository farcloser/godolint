package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3025 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3025Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3025(t *testing.T) {
	allRules := []rule.Rule{DL3025()}

	t.Run("don't warn on CMD JSON notation with broken long strings", func(t *testing.T) {
		dockerfile := `CMD [ "/bin/sh", "-c", \
      "echo foo && \
       echo bar" \
    ]
DL3025`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3025")
	})

	t.Run("don't warn on CMD json notation", func(t *testing.T) {
		dockerfile := `FROM scratch as build
CMD ["foo", "bar"]
CMD [ "foo", "bar" ]
DL3025`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3025")
	})

	t.Run("don't warn on ENTRYPOINT json notation", func(t *testing.T) {
		dockerfile := `FROM scratch as build
ENTRYPOINT ["foo", "bar"]
DL3025`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3025")
	})

	t.Run("warn on CMD", func(t *testing.T) {
		dockerfile := `FROM node as foo
CMD something
DL3025`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3025")
	})

	t.Run("warn on ENTRYPOINT", func(t *testing.T) {
		dockerfile := `FROM node as foo
ENTRYPOINT something
DL3025`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3025")
	})
}
