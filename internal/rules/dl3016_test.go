package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3016 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3016Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3016(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3016(),
	}

	t.Run(
		"commit not pinned for git",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install git://github.com/npm/npm.git`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"commit not pinned for git+http",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install git+http://isaacs@github.com/npm/npm`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"commit not pinned for git+ssh",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install git+ssh://git@github.com:npm/npm.git`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"don't fire on loglevel flag",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install --loglevel verbose sax@0.1.1`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for folder - absolute path",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install /folder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for folder - relative path from current folder",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install ./folder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for folder - relative path from home",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install ~/folder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for folder - relative path to parent folder",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install ../folder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for tarball suffix .tar",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install package-v1.2.3.tar`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for tarball suffix .tar.gz",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install package-v1.2.3.tar.gz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version does not have to be pinned for tarball suffix .tgz",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install package-v1.2.3.tgz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version not pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install express`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version not pinned multiple packages",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install express sax@0.1.1`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version not pinned with --global",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install --global express`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version not pinned with scope",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install @myorg/privatepackage`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install express@4.1.1`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned in package.json",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned in package.json with arguments",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install --progress=false`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned multiple packages",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install express@"4.1.1" sax@0.1.1`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned with --global",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install --global express@"4.1.1"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned with -g",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install -g express@"4.1.1"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned with scope",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install @myorg/privatepackage@">=0.1.0 <0.2.0"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)

	t.Run(
		"version pinned with scope",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `RUN npm install @myorg/privatepackage@">=0.1.0"`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3016")
		},
	)
}
