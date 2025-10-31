package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3014 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3014Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3014(t *testing.T) {
	allRules := []rule.Rule{DL3014()}

	t.Run("apt-get --quiet", func(t *testing.T) {
		dockerfile := `RUN apt-get install --quiet python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3014")
	})

	t.Run("apt-get --quiet --quiet", func(t *testing.T) {
		dockerfile := `RUN apt-get install --quiet --quiet python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get -q", func(t *testing.T) {
		dockerfile := `RUN apt-get install -q python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3014")
	})

	t.Run("apt-get -q -q", func(t *testing.T) {
		dockerfile := `RUN apt-get install -q -q python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get -q=2", func(t *testing.T) {
		dockerfile := `RUN apt-get install -q=2 python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get -qq", func(t *testing.T) {
		dockerfile := `RUN apt-get install -qq python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get auto yes", func(t *testing.T) {
		dockerfile := `RUN apt-get install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3014")
	})

	t.Run("apt-get with assume-yes", func(t *testing.T) {
		dockerfile := `RUN apt-get --assume-yes install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get with auto expanded yes", func(t *testing.T) {
		dockerfile := `RUN apt-get --yes install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get with auto yes", func(t *testing.T) {
		dockerfile := `RUN apt-get -y install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get yes different pos", func(t *testing.T) {
		dockerfile := `RUN apt-get install -y python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})

	t.Run("apt-get yes shortflag", func(t *testing.T) {
		dockerfile := `RUN apt-get install -yq python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3014")
	})
}
