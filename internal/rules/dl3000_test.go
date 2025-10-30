package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3000 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3000Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3000(t *testing.T) {
	allRules := []rule.Rule{ DL3000() }


	t.Run("workdir absolute", func(t *testing.T) {
		dockerfile := `WORKDIR /usr/local`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir absolute double quotes", func(t *testing.T) {
		dockerfile := `WORKDIR "/usr/local"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir absolute single quotes", func(t *testing.T) {
		dockerfile := `WORKDIR '/usr/local'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir absolute windows", func(t *testing.T) {
		dockerfile := `WORKDIR 'C:\'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir absolute windows alternative", func(t *testing.T) {
		dockerfile := `WORKDIR C:/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir absolute windows quotes", func(t *testing.T) {
		dockerfile := `WORKDIR "C:\"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir absolute windows quotes alternative", func(t *testing.T) {
		dockerfile := `WORKDIR "C:/"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir relative", func(t *testing.T) {
		dockerfile := `WORKDIR relative/dir`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3000")

	})

	t.Run("workdir relative double quotes", func(t *testing.T) {
		dockerfile := `WORKDIR "relative/dir"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3000")

	})

	t.Run("workdir relative single quotes", func(t *testing.T) {
		dockerfile := `WORKDIR 'relative/dir'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3000")

	})

	t.Run("workdir variable", func(t *testing.T) {
		dockerfile := `WORKDIR ${work}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

	t.Run("workdir variable double quotes", func(t *testing.T) {
		dockerfile := `WORKDIR "${dir}"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3000")

	})

}
