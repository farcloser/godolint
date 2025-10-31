package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3009 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3009Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3009(t *testing.T) {
	allRules := []rule.Rule{ DL3009() }


	t.Run("apt cleanup", func(t *testing.T) {
		dockerfile := `apt cleanup
FROM scratch
RUN apt update && apt install python && rm -rf /var/lib/apt/lists/*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("apt no cleanup", func(t *testing.T) {
		dockerfile := `apt no cleanup
FROM scratch
RUN apt update && apt install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3009")

	})

	t.Run("apt-get cleanup", func(t *testing.T) {
		dockerfile := `apt-get cleanup
FROM scratch
RUN apt-get update && apt-get install python && rm -rf /var/lib/apt/lists/*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("apt-get cleanup in stage image", func(t *testing.T) {
		dockerfile := `apt-get cleanup in stage image
FROM ubuntu as foo
RUN apt-get update && apt-get install python
FROM scratch
RUN echo hey!`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("apt-get no cleanup", func(t *testing.T) {
		dockerfile := `apt-get no cleanup
FROM scratch
RUN apt-get update && apt-get install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3009")

	})

	t.Run("apt-get no cleanup in intermediate stage", func(t *testing.T) {
		dockerfile := `apt-get no cleanup in intermediate stage
FROM ubuntu as foo
RUN apt-get update && apt-get install python
FROM foo
RUN hey!`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3009")

	})

	t.Run("apt-get no cleanup in last stage", func(t *testing.T) {
		dockerfile := `apt-get no cleanup in last stage
FROM ubuntu as foo
RUN hey!
FROM scratch
RUN apt-get update && apt-get install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3009")

	})

	t.Run("aptitude cleanup", func(t *testing.T) {
		dockerfile := `aptitude cleanup
FROM scratch
RUN aptitude update && aptitude install python && rm -rf /var/lib/apt/lists/*`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("aptitude no cleanup", func(t *testing.T) {
		dockerfile := `aptitude no cleanup
FROM scratch
RUN aptitude update && aptitude install python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3009")

	})

	t.Run("don't warn: cache mount to apt cache and lists directory", func(t *testing.T) {
		dockerfile := `don't warn: cache mount to apt cache and lists directory
RUN \
  --mount=type=cache,target=/var/cache/apt \
  --mount=type=cache,target=/var/lib/apt \
  apt-get update`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("don't warn: cache mount to apt cache and tmpfs mount to lists directory", func(t *testing.T) {
		dockerfile := `don't warn: cache mount to apt cache and tmpfs mount to lists directory
RUN \
  --mount=type=cache,target=/var/cache/apt \
  --mount=type=tmpfs,target=/var/lib/apt \
  apt-get update`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("don't warn: tmpfs mount to apt cache and cache mount to lists directory", func(t *testing.T) {
		dockerfile := `don't warn: tmpfs mount to apt cache and cache mount to lists directory
RUN \
  --mount=type=tmpfs,target=/var/cache/apt \
  --mount=type=cache,target=/var/lib/apt \
  apt-get update`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("don't warn: tmpfs mount to apt cache and lists directory", func(t *testing.T) {
		dockerfile := `don't warn: tmpfs mount to apt cache and lists directory
RUN \
  --mount=type=tmpfs,target=/var/cache/apt \
  --mount=type=tmpfs,target=/var/lib/apt \
  apt-get update`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("no warn apt-get cleanup in intermediate stage that cleans lists", func(t *testing.T) {
		dockerfile := `no warn apt-get cleanup in intermediate stage that cleans lists
FROM ubuntu as foo
RUN apt-get update && apt-get install python && rm -rf /var/lib/apt/lists/*
FROM foo
RUN hey!`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

	t.Run("no warn apt-get cleanup in intermediate stage when stage not used later", func(t *testing.T) {
		dockerfile := `no warn apt-get cleanup in intermediate stage when stage not used later
FROM ubuntu as foo
RUN apt-get update && apt-get install python
FROM scratch
RUN hey!`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3009")

	})

}
