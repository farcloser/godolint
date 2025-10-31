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


	t.Run("apt-get pinned chained", func(t *testing.T) {
		dockerfile := `apt-get pinned chained
RUN apt-get update \
 && apt-get -yqq --no-install-recommends install nodejs=0.10 \
 && rm -rf /var/lib/apt/lists/*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3008")

	})

	t.Run("apt-get pinned regression", func(t *testing.T) {
		dockerfile := `apt-get pinned regression
RUN apt-get update && apt-get install --no-install-recommends -y \
python-demjson=2.2.2* \
wget=1.16.1* \
git=1:2.5.0* \
ruby=1:2.1.*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3008")

	})

	t.Run("apt-get tolerate target-release", func(t *testing.T) {
		dockerfile := `apt-get tolerate target-release
RUN set -e &&\
 apt-get update &&\
 rm -rf /var/lib/apt/lists/*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3008")

	})

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
