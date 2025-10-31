package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3013 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3013Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3013(t *testing.T) {
	allRules := []rule.Rule{DL3013()}

	t.Run("pip install constraints file - long version argument", func(t *testing.T) {
		dockerfile := `RUN pip install pykafka --constraint http://foo.bar.baz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install constraints file - short version argument", func(t *testing.T) {
		dockerfile := `RUN pip install pykafka -c http://foo.bar.baz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install excluded version", func(t *testing.T) {
		dockerfile := `RUN pip install 'alabaster!=0.7'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install lower bound", func(t *testing.T) {
		dockerfile := `RUN pip install 'alabaster<0.7'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install no cache dir", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python==1.2.2 --no-cache-dir`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install requirements", func(t *testing.T) {
		dockerfile := `RUN pip install -r requirements.txt`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install requirements with long flag", func(t *testing.T) {
		dockerfile := `RUN pip install --requirement requirements.txt`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install upper bound", func(t *testing.T) {
		dockerfile := `RUN pip install 'alabaster>=0.7'`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install use setup.py", func(t *testing.T) {
		dockerfile := `RUN pip install .`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip install user directory", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python==1.2.2 --user`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version not pinned", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3013")
	})

	t.Run("pip version not pinned with python -m", func(t *testing.T) {
		dockerfile := `RUN python -m pip install example`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python==1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with === operator", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python===1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with flag --build", func(t *testing.T) {
		dockerfile := `RUN pip3 install --build /opt/yamllint yamllint==1.20.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with flag --ignore-installed", func(t *testing.T) {
		dockerfile := `RUN pip install --ignore-installed MySQL_python==1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with flag --prefix", func(t *testing.T) {
		dockerfile := `RUN pip3 install --prefix /opt/yamllint yamllint==1.20.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with flag --root", func(t *testing.T) {
		dockerfile := `RUN pip3 install --root /opt/yamllint yamllint==1.20.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with flag --target", func(t *testing.T) {
		dockerfile := `RUN pip3 install --target /opt/yamllint yamllint==1.20.0`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with flag --trusted-host", func(t *testing.T) {
		dockerfile := `RUN pip3 install --trusted-host host example==1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with python -m", func(t *testing.T) {
		dockerfile := `RUN python -m pip install example==1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip version pinned with ~= operator", func(t *testing.T) {
		dockerfile := `RUN pip install MySQL_python~=1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip2 version not pinned", func(t *testing.T) {
		dockerfile := `RUN pip2 install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3013")
	})

	t.Run("pip3 install from local package", func(t *testing.T) {
		dockerfile := `RUN pip3 install mypkg.whl`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip3 install from local package (2)", func(t *testing.T) {
		dockerfile := `RUN pip3 install mypkg.tar.gz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pip3 version not pinned", func(t *testing.T) {
		dockerfile := `RUN pip3 install MySQL_python`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3013")
	})

	t.Run("pip3 version pinned", func(t *testing.T) {
		dockerfile := `RUN pip3 install MySQL_python==1.2.2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})

	t.Run("pipenv is not pip", func(t *testing.T) {
		dockerfile := `RUN pipenv install black`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3013")
	})
}
