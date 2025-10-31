package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3030 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3030Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3030(t *testing.T) {
	allRules := []rule.Rule{DL3030()}

	t.Run("not ok when not using `-y` switch", func(t *testing.T) {
		dockerfile := `RUN yum install httpd-2.4.24 && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3030")
	})

	t.Run("ok when using `-y` switch", func(t *testing.T) {
		dockerfile := `RUN yum install -y httpd-2.4.24 && yum clean all`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3030")
	})

	t.Run("ok when using `-y` switch (2)", func(t *testing.T) {
		dockerfile := `RUN bash -c ` + "`" + `# not even a yum command` + "`" + ``
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3030")
	})
}
