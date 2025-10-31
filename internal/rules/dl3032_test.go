package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3032 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3032Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3032(t *testing.T) {
	allRules := []rule.Rule{DL3032()}

	t.Run("not ok with no clean all", func(t *testing.T) {
		dockerfile := `RUN yum install -y mariadb-10.4`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3032")
	})

	t.Run("ok with rm -rf /var/cache/yum/*", func(t *testing.T) {
		dockerfile := `RUN yum install -y mariadb-10.4 && rm -rf /var/cache/yum/*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3032")
	})

	t.Run("ok with yum clean all ", func(t *testing.T) {
		dockerfile := `RUN yum install -y mariadb-10.4 && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3032")
	})

	t.Run("ok with yum clean all  (2)", func(t *testing.T) {
		dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3032")
	})
}
