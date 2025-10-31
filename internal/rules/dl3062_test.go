package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3062 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3062Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3062(t *testing.T) {
	allRules := []rule.Rule{DL3062()}

	t.Run("go version not pinned", func(t *testing.T) {
		dockerfile := `RUN go install example.com/pkg`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3062")
	})

	t.Run("go version not pinned", func(t *testing.T) {
		dockerfile := `RUN go get example.com/pkg`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3062")
	})

	t.Run("go version not pinned", func(t *testing.T) {
		dockerfile := `RUN go run example.com/pkg`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3062")
	})

	t.Run("go version pinned", func(t *testing.T) {
		dockerfile := `RUN go install example.com/pkg@v1.2.3`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3062")
	})

	t.Run("go version pinned", func(t *testing.T) {
		dockerfile := `RUN go get example.com/pkg@v1.2.3`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3062")
	})

	t.Run("go version pinned", func(t *testing.T) {
		dockerfile := `RUN go run example.com/pkg@v1.2.3`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3062")
	})

	t.Run("go version pinned as latest", func(t *testing.T) {
		dockerfile := `RUN go install example.com/pkg@latest`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3062")
	})

	t.Run("go version pinned as latest", func(t *testing.T) {
		dockerfile := `RUN go get example.com/pkg@latest`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3062")
	})

	t.Run("go version pinned as latest", func(t *testing.T) {
		dockerfile := `RUN go run example.com/pkg@latest`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3062")
	})
}
