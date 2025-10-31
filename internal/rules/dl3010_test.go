package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3010 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3010Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3010(t *testing.T) {
	allRules := []rule.Rule{DL3010()}

	t.Run("catch: copy archive then extract 1", func(t *testing.T) {
		dockerfile := `COPY packaged-app.tar /usr/src/app
RUN tar -xf /usr/src/app/packaged-app.tar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3010")
	})

	t.Run("catch: copy archive then extract 2", func(t *testing.T) {
		dockerfile := `COPY packaged-app.tar /usr/src/app
WORKDIR /usr/src/app
RUN foo bar && echo something && tar -xf packaged-app.tar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3010")
	})

	t.Run("catch: copy archive then extract 3", func(t *testing.T) {
		dockerfile := `COPY foo/bar/packaged-app.tar /foo.tar
RUN tar -xf /foo.tar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3010")
	})

	t.Run("catch: copy archive then extract windows paths 1", func(t *testing.T) {
		dockerfile := `COPY build\foo\bar.tar.gz "C:\Program Files\Foo"
RUN tar -xf "C:\Program Files\Foo\bar.tar.gz"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3010")
	})

	t.Run("catch: copy archive then extract windows paths 2", func(t *testing.T) {
		dockerfile := `COPY build\foo\bar.tar.gz "C:\Program Files\foo.tar.gz"
RUN tar -xf foo.tar.gz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3010")
	})

	t.Run("ignore: copy archive without extract", func(t *testing.T) {
		dockerfile := `COPY packaged-app.tar /usr/src/app
FROM debian:11 as newstage`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3010")
	})

	t.Run("ignore: copy from previous stage", func(t *testing.T) {
		dockerfile := `COPY --from=builder /usr/local/share/some.tar /opt/some.tar`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3010")
	})

	t.Run("ignore: non archive", func(t *testing.T) {
		dockerfile := `COPY package.json /usr/src/app`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3010")
	})
}
