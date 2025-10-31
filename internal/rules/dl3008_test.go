package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3008 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3008Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3008(t *testing.T) {
	allRules := []rule.Rule{ DL3008() }


	t.Run("apt-get version", func(t *testing.T) {
		dockerfile := `RUN apt-get install -y python=1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3008")

	})

	t.Run("apt-get version", func(t *testing.T) {
		dockerfile := `RUN apt-get install ./wkhtmltox_0.12.5-1.bionic_amd64.deb`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3008")

	})

	t.Run("apt-get version pinning", func(t *testing.T) {
		dockerfile := `RUN apt-get update && apt-get install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3008")

	})

}
