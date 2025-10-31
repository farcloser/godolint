package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3037 checks for zypper install without version pinning.
func DL3037() rule.Rule {
	return rule.NewSimpleRule(
		DL3037Meta.Code,
		DL3037Meta.Severity,
		DL3037Meta.Message,
		checkDL3037,
	)
}

func checkDL3037(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all zypper packages
	packages := getZypperPackages(parsed)
	for _, pkg := range packages {
		if !isZypperVersionFixed(pkg) {
			return false
		}
	}

	return true
}

func getZypperPackages(parsed *shell.ParsedShell) []string {
	var packages []string

	for _, cmd := range parsed.PresentCommands {
		if shell.CmdHasArgs("zypper", []string{"install"}, cmd) ||
			shell.CmdHasArgs("zypper", []string{"in"}, cmd) {
			args := shell.GetArgsNoFlags(cmd)
			for _, arg := range args {
				if arg != "install" && arg != "in" {
					packages = append(packages, arg)
				}
			}
		}
	}

	return packages
}

func isZypperVersionFixed(pkg string) bool {
	// Check for version operators
	if strings.Contains(pkg, "=") ||
		strings.Contains(pkg, ">=") ||
		strings.Contains(pkg, ">") ||
		strings.Contains(pkg, "<=") ||
		strings.Contains(pkg, "<") {
		return true
	}

	// .rpm files are considered version-fixed
	if strings.HasSuffix(pkg, ".rpm") {
		return true
	}

	return false
}
