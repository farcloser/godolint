package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3043 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3043Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3043(t *testing.T) {
	allRules := []rule.Rule{ DL3043() }


	t.Run("error when using `FROM` within `ONBUILD`", func(t *testing.T) {
		dockerfile := `error when using ` + "`" + `FROM` + "`" + ` within ` + "`" + `ONBUILD` + "`" + `
ONBUILD FROM debian:buster
DL3043`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3043")

	})

	t.Run("error when using `MAINTAINER` within `ONBUILD`", func(t *testing.T) {
		dockerfile := `error when using ` + "`" + `MAINTAINER` + "`" + ` within ` + "`" + `ONBUILD` + "`" + `
ONBUILD MAINTAINER "BoJack Horseman"
DL3043`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3043")

	})

	t.Run("error when using `ONBUILD` within `ONBUILD`", func(t *testing.T) {
		dockerfile := `error when using ` + "`" + `ONBUILD` + "`" + ` within ` + "`" + `ONBUILD` + "`" + `
ONBUILD ONBUILD RUN anything
DL3043`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3043")

	})

	t.Run("ok with `ADD`", func(t *testing.T) {
		dockerfile := `ONBUILD ADD anything anywhere`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `ARG`", func(t *testing.T) {
		dockerfile := `ONBUILD ARG anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `CMD`", func(t *testing.T) {
		dockerfile := `ONBUILD CMD anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `COPY`", func(t *testing.T) {
		dockerfile := `ONBUILD COPY anything anywhere`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `ENTRYPOINT`", func(t *testing.T) {
		dockerfile := `ONBUILD ENTRYPOINT anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `ENV`", func(t *testing.T) {
		dockerfile := `ONBUILD ENV MYVAR="bla"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `EXPOSE`", func(t *testing.T) {
		dockerfile := `ONBUILD EXPOSE 69`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `FROM` outside of `ONBUILD`", func(t *testing.T) {
		dockerfile := `FROM debian:buster`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `HEALTHCHECK`", func(t *testing.T) {
		dockerfile := `ONBUILD HEALTHCHECK NONE`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `LABEL`", func(t *testing.T) {
		dockerfile := `ONBUILD LABEL bla="blubb"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `MAINTAINER` outside of `ONBUILD`", func(t *testing.T) {
		dockerfile := `MAINTAINER "Some Guy"`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `RUN`", func(t *testing.T) {
		dockerfile := `ONBUILD RUN anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `SHELL`", func(t *testing.T) {
		dockerfile := `ONBUILD SHELL anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `STOPSIGNAL`", func(t *testing.T) {
		dockerfile := `ONBUILD STOPSIGNAL anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `USER`", func(t *testing.T) {
		dockerfile := `ONBUILD USER anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `VOLUME`", func(t *testing.T) {
		dockerfile := `ONBUILD VOLUME anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

	t.Run("ok with `WORKDIR`", func(t *testing.T) {
		dockerfile := `ONBUILD WORKDIR anything`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3043")

	})

}
