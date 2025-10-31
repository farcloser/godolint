package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3036 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3036Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3036(t *testing.T) {
	allRules := []rule.Rule{ DL3036() }


	t.Run("not ok without zypper clean", func(t *testing.T) {
		dockerfile := `RUN zypper install -y mariadb=10.4`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3036")

	})

	t.Run("ok when mount type cache is used", func(t *testing.T) {
		dockerfile := `RUN --mount=type=cache,target=/var/cache/zypp zypper install -y mariadb`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3036")

	})

	t.Run("ok when mount type tmpfs is used", func(t *testing.T) {
		dockerfile := `RUN --mount=type=tmpfs,target=/var/cache/zypp zypper install -y mariadb`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3036")

	})

	t.Run("ok with zypper clean", func(t *testing.T) {
		dockerfile := `RUN zypper install -y mariadb=10.4 && zypper clean`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3036")

	})

	t.Run("ok with zypper clean (2)", func(t *testing.T) {
		dockerfile := `RUN zypper install -y mariadb=10.4 && zypper cc`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3036")

	})

}
