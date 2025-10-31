package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3042 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3042Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3042(t *testing.T) {
	allRules := []rule.Rule{ DL3042() }


	t.Run("don't match on pipenv", func(t *testing.T) {
		dockerfile := `RUN pipenv install library`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("don't match on pipx", func(t *testing.T) {
		dockerfile := `RUN pipx install software`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("ok with cache mount in cache dir", func(t *testing.T) {
		dockerfile := `RUN --mount=type=cache,target=/root/.cache/pip pip install foobar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("ok with tmpfs mount in cache dir", func(t *testing.T) {
		dockerfile := `RUN --mount=type=tmpfs,target=/root/.cache/pip pip install foobar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("pip --no-cache-dir not used", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("pip --no-cache-dir used", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python --no-cache-dir`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("pip2 --no-cache-dir not used", func(t *testing.T) {
		dockerfile := `RUN pip2 install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("pip2 --no-cache-dir used", func(t *testing.T) {
		dockerfile := `RUN pip2 install MySQL_python --no-cache-dir`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("pip3 --no-cache-dir not used", func(t *testing.T) {
		dockerfile := `RUN pip3 install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("pip3 --no-cache-dir used", func(t *testing.T) {
		dockerfile := `RUN pip3 install --no-cache-dir MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with falsy values", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=0
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with falsy values (2)", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=off
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with falsy values (3)", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=no
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with falsy values (4)", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=false
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with truthy values", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=1
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with truthy values (2)", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=on
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with truthy values (3)", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=yes
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect ENV PIP_NO_CACHE_DIR with truthy values (4)", func(t *testing.T) {
		dockerfile := `ENV PIP_NO_CACHE_DIR=true
RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with falsy values", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=0 pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with falsy values (2)", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=off pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with falsy values (3)", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=no pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with falsy values (4)", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=false pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with truthy values", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=1 pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with truthy values (2)", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=on pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with truthy values (3)", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=yes pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN PIP_NO_CACHE_DIR=... with truthy values (4)", func(t *testing.T) {
		dockerfile := `RUN PIP_NO_CACHE_DIR=true pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with falsy values", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=0 && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with falsy values (2)", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=off && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with falsy values (3)", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=no && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with falsy values (4)", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=false && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with truthy values", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=1 && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with truthy values (2)", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=on && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with truthy values (3)", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=yes && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

	t.Run("respect RUN export PIP_NO_CACHE_DIR=... with truthy values (4)", func(t *testing.T) {
		dockerfile := `RUN export PIP_NO_CACHE_DIR=true && pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3042")

	})

}
