package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3048 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3048Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3048(t *testing.T) {
	allRules := []rule.Rule{DL3048()}

	t.Run("not ok with consecutive dividers", func(t *testing.T) {
		dockerfile := `LABEL invalid..character="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with consecutive dividers (2)", func(t *testing.T) {
		dockerfile := `LABEL invalid--character="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with invalid character", func(t *testing.T) {
		dockerfile := `LABEL invalid$character="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with invalid start and end characters", func(t *testing.T) {
		dockerfile := `LABEL .invalid ="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with invalid start and end characters (2)", func(t *testing.T) {
		dockerfile := `LABEL -invalid ="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with invalid start and end characters (3)", func(t *testing.T) {
		dockerfile := `LABEL 1invalid ="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with reserved namespace", func(t *testing.T) {
		dockerfile := `LABEL com.docker.label="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with reserved namespace (2)", func(t *testing.T) {
		dockerfile := `LABEL io.docker.label="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("not ok with reserved namespace (3)", func(t *testing.T) {
		dockerfile := `LABEL org.dockerproject.label="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3048")
	})

	t.Run("ok with valid labels", func(t *testing.T) {
		dockerfile := `LABEL org.valid-key.label3="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3048")
	})

	t.Run("ok with valid labels (2)", func(t *testing.T) {
		dockerfile := `LABEL validlabel="foo"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3048")
	})
}
