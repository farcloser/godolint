package rules

import (
	"strings"

	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/shell"
	"github.com/forkcloser/godolint/internal/syntax"
)

// DL3062 checks for go install, go get and go run without a pinned version.
func DL3062() rule.Rule {
	return rule.NewSimpleRule(
		DL3062Meta.Code,
		DL3062Meta.Severity,
		DL3062Meta.Message,
		checkDL3062,
	)
}

func checkDL3062(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	for _, pkg := range goPackages(parsed) {
		if !isGoVersionPinned(pkg) {
			return false
		}
	}

	return true
}

// goPackages is what hadolint checks: the package a `go run` names when
// there is one, else the packages `go install` and `go get` name.
func goPackages(parsed *shell.ParsedShell) []string {
	if run := goRunPackages(parsed); len(run) > 0 {
		return run
	}

	return goInstallPackages(parsed)
}

// goRunPackages is the first non-flag argument after `run` of each `go run`:
// everything after the package is an argument to the program it runs.
func goRunPackages(parsed *shell.ParsedShell) []string {
	var packages []string

	for _, cmd := range parsed.PresentCommands {
		if !shell.CmdHasArgs("go", []string{"run"}, cmd) {
			continue
		}

		for idx, arg := range shell.GetArgsNoFlags(cmd) {
			if arg != "run" && idx <= 1 {
				packages = append(packages, arg)
			}
		}
	}

	return packages
}

// goInstallPackages is every non-flag argument of `go install` and `go get`
// but the subcommand words.
func goInstallPackages(parsed *shell.ParsedShell) []string {
	var packages []string

	for _, cmd := range parsed.PresentCommands {
		if !shell.CmdHasArgs("go", []string{installArg, "get"}, cmd) {
			continue
		}

		for _, arg := range shell.GetArgsNoFlags(cmd) {
			if arg != installArg && arg != "get" && arg != "tool" {
				packages = append(packages, arg)
			}
		}
	}

	return packages
}

// isGoVersionPinned is a local path, which has no version to pin, or a module
// path with an @version that is neither @latest nor @none.
func isGoVersionPinned(pkg string) bool {
	if isGoLocalPath(pkg) {
		return true
	}

	if !strings.Contains(pkg, "@") {
		return false
	}

	return !strings.HasSuffix(pkg, "@latest") && !strings.HasSuffix(pkg, "@none")
}

// isGoLocalPath is `.`, or a path that starts with `/` or `.`.
func isGoLocalPath(pkg string) bool {
	return pkg == "." || strings.HasPrefix(pkg, "/") || strings.HasPrefix(pkg, ".")
}
