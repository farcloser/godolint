package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3062 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3062Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3062(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3062(),
	}

	t.Run(
		"go version not pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go install example.com/pkg`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version not pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go get example.com/pkg`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version not pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go run example.com/pkg`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go install example.com/pkg@v1.2.3`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go get example.com/pkg@v1.2.3`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go run example.com/pkg@v1.2.3`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version pinned as latest",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go install example.com/pkg@latest`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version pinned as latest",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go get example.com/pkg@latest`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3062")
		},
	)

	t.Run(
		"go version pinned as latest",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN go run example.com/pkg@latest`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3062")
		},
	)
}
