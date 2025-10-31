package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3045 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3045Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3045(t *testing.T) {
	allRules := []rule.Rule{DL3045()}

	t.Run("not ok: `COPY` with relative destination and no `WORKDIR` set", func(t *testing.T) {
		dockerfile := `COPY bla.sh blubb.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3045")
	})

	t.Run("not ok: `COPY` with relative destination and no `WORKDIR` set with quotes", func(t *testing.T) {
		dockerfile := `COPY bla.sh "blubb.sh"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with absolute destination and no `WORKDIR` set", func(t *testing.T) {
		dockerfile := `COPY bla.sh /usr/local/bin/blubb.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with absolute destination and no `WORKDIR` set - windows", func(t *testing.T) {
		dockerfile := `COPY bla.sh c:\system32\blubb.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run(
		"ok: `COPY` with absolute destination and no `WORKDIR` set - windows with alternative paths",
		func(t *testing.T) {
			dockerfile := `COPY bla.sh c:/system32/blubb.sh`
			violations := LintDockerfile(dockerfile, allRules)

			AssertNoViolation(t, violations, "DL3045")
		},
	)

	t.Run("ok: `COPY` with absolute destination and no `WORKDIR` set - windows with quotes", func(t *testing.T) {
		dockerfile := `COPY bla.sh "c:\system32\blubb.sh"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with absolute destination and no `WORKDIR` set with quotes", func(t *testing.T) {
		dockerfile := `COPY bla.sh "/usr/local/bin/blubb.s"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with destination being an environment variable 1", func(t *testing.T) {
		dockerfile := `COPY src.sh ${SRC_BASE_ENV}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with destination being an environment variable 2", func(t *testing.T) {
		dockerfile := `COPY src.sh $SRC_BASE_ENV`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with destination being an environment variable 3", func(t *testing.T) {
		dockerfile := `COPY src.sh "${SRC_BASE_ENV}"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with destination being an environment variable 4", func(t *testing.T) {
		dockerfile := `COPY src.sh "$SRC_BASE_ENV"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with relative destination and `WORKDIR` set", func(t *testing.T) {
		dockerfile := `FROM scratch
WORKDIR /usr
COPY bla.sh blubb.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("ok: `COPY` with relative destination and `WORKDIR` set - windows", func(t *testing.T) {
		dockerfile := `FROM scratch
WORKDIR c:\system32
COPY bla.sh blubb.sh`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3045")
	})

	t.Run("regression: don't crash with single character paths", func(t *testing.T) {
		dockerfile := `COPY a b`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3045")
	})
}
