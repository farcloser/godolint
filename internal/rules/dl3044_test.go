package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3044 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3044Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3044(t *testing.T) {
	allRules := []rule.Rule{DL3044()}

	t.Run("fail when referencing a variable on its own right side twice within the same `ENV`", func(t *testing.T) {
		dockerfile := `ENV PATH=/bla:${PATH} PATH=/blubb:${PATH}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3044")
	})

	t.Run("fail with full match 1", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="$BLA"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3044")
	})

	t.Run("fail with full match 2", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="${BLA}"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3044")
	})

	t.Run("fail with partial match 5", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="$BLA/$BLAFOO/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3044")
	})

	t.Run("fail with selfreferencing with curly braces ENV", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="${BLA}/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3044")
	})

	t.Run("fail with selfreferencing without curly braces ENV", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="$BLA/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3044")
	})

	t.Run("ok when previously defined in `ARG`", func(t *testing.T) {
		dockerfile := `ARG BLA
ENV BLA=${BLA}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok when previously defined in `ENV`", func(t *testing.T) {
		dockerfile := `ENV BLA blubb
ENV BLA=${BLA}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with normal ENV", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb"
ENV BLUBB="${BLA}/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with parial match 6", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="BLA/$BLAFOO/BLA"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with partial match 1", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="${FOOBLA}/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with partial match 2", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="${BLAFOO}/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with partial match 3", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="$FOOBLA/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with partial match 4", func(t *testing.T) {
		dockerfile := `ENV BLA="blubb" BLUBB="$BLAFOO/blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with referencing a variable on its own right hand side", func(t *testing.T) {
		dockerfile := `ENV PATH=/bla:${PATH}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})

	t.Run("ok with referencing a variable on its own right side twice in different `ENV`s", func(t *testing.T) {
		dockerfile := `ENV PATH=/bla:${PATH}
ENV PATH=/blubb:${PATH}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3044")
	})
}
