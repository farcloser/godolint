package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3020 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3020Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3020(t *testing.T) {
	allRules := []rule.Rule{ DL3020() }


	t.Run("add for bz2", func(t *testing.T) {
		dockerfile := `ADD file.bz2 /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for gzip", func(t *testing.T) {
		dockerfile := `ADD file.gz /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for tar", func(t *testing.T) {
		dockerfile := `ADD file.tar /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for tgz", func(t *testing.T) {
		dockerfile := `ADD file.tgz /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for tgz with quotes", func(t *testing.T) {
		dockerfile := `ADD "file.tgz" /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for url", func(t *testing.T) {
		dockerfile := `ADD http://file.com /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for url with quotes", func(t *testing.T) {
		dockerfile := `ADD "http://file.com" /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("add for xz", func(t *testing.T) {
		dockerfile := `ADD file.xz /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3020")

	})

	t.Run("using add", func(t *testing.T) {
		dockerfile := `ADD file /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3020")

	})

	t.Run("warn for zip", func(t *testing.T) {
		dockerfile := `ADD file.zip /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3020")

	})

	t.Run("warn for zip with quotes", func(t *testing.T) {
		dockerfile := `ADD "file.zip" /usr/src/app/`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3020")

	})

}
