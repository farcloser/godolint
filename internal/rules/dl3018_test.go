package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3018 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3018Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3018(t *testing.T) {
	allRules := []rule.Rule{DL3018()}

	t.Run("apk add no version pinning single", func(t *testing.T) {
		dockerfile := `RUN apk add flex=2.6.4-r1`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("apk add version pinned chained", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache flex=2.6.4-r1 \
 && pip install -r requirements.txt`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("apk add version pinned regression", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache \
flex=2.6.4-r1 \
libffi=3.2.1-r3 \
python2=2.7.13-r1 \
libbz2=1.0.6-r5`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("apk add version pinned regression - one missed", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache \
flex=2.6.4-r1 \
libffi \
python2=2.7.13-r1 \
libbz2=1.0.6-r5`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3018")
	})

	t.Run("apk add version pinning single", func(t *testing.T) {
		dockerfile := `RUN apk add flex`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3018")
	})

	t.Run("apk add virtual package", func(t *testing.T) {
		dockerfile := `RUN apk add \
--virtual build-dependencies \
python-dev=1.1.1 build-base=2.2.2 wget=3.3.3 \
&& pip install -r requirements.txt \
&& python setup.py install \
&& apk del build-dependencies`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("apk add with repository (-X) without equal sign", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache \
-X https://nl.alpinelinux.org/alpine/edge/testing \
flow=0.78.0-r0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("apk add with repository with equal sign", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache \
--repository=https://nl.alpinelinux.org/alpine/edge/testing \
flow=0.78.0-r0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("apk add with repository without equal sign", func(t *testing.T) {
		dockerfile := `RUN apk add --no-cache \
--repository https://nl.alpinelinux.org/alpine/edge/testing \
flow=0.78.0-r0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})

	t.Run("don't trigger when installing from .apk file", func(t *testing.T) {
		dockerfile := `RUN apk add mypackage-1.1.1.apk`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3018")
	})
}
