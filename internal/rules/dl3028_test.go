package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3028 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3028Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3028(t *testing.T) {
	allRules := []rule.Rule{ DL3028() }


	t.Run("does not warn on --version with =", func(t *testing.T) {
		dockerfile := `RUN gem install bundler --version='2.0.1'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("does not warn on --version without =", func(t *testing.T) {
		dockerfile := `RUN gem install bundler --version '2.0.1'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("does not warn on -v", func(t *testing.T) {
		dockerfile := `RUN gem install bundler -v '2.0.1'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("does not warn when using extra flags", func(t *testing.T) {
		dockerfile := `RUN gem install bundler:2.0.1 --use-system-libraries true`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("does not warn when using extra flags with double dashes", func(t *testing.T) {
		dockerfile := `RUN gem install bundler:2.0.1 -- --use-system-libraries true`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("does not warn when using extra flags with equal sign", func(t *testing.T) {
		dockerfile := `RUN gem install bundler:2.0.1 --use-system-libraries=true`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("does not warn when using extra flags with equal sign and double dashes", func(t *testing.T) {
		dockerfile := `RUN gem install bundler:2.0.1 -- --use-system-libraries=true`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("multi", func(t *testing.T) {
		dockerfile := `RUN gem i bunlder:1 nokogiri`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3028")

	})

	t.Run("multi (2)", func(t *testing.T) {
		dockerfile := `RUN gem i bunlder:1 nokogirii:1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("pinned", func(t *testing.T) {
		dockerfile := `RUN gem install bundler:1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("pinned", func(t *testing.T) {
		dockerfile := `RUN gem i bundler:1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3028")

	})

	t.Run("unpinned", func(t *testing.T) {
		dockerfile := `RUN gem install bundler`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3028")

	})

	t.Run("unpinned", func(t *testing.T) {
		dockerfile := `RUN gem i bundler`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3028")

	})

}
