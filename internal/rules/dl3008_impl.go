package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3008 checks for apt-get install without version pinning.
func DL3008() rule.Rule {
	return rule.NewSimpleRule(
		DL3008Meta.Code,
		DL3008Meta.Severity,
		DL3008Meta.Message,
		checkDL3008,
	)
}

func checkDL3008(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all apt-get install commands
	for _, cmd := range parsed.PresentCommands {
		packages := aptGetPackages(cmd)
		for _, pkg := range packages {
			if !versionFixed(pkg) {
				return false
			}
		}
	}

	return true
}

func aptGetPackages(cmd shell.Command) []string {
	// Only check apt-get install commands
	if !shell.CmdHasArgs("apt-get", []string{"install"}, cmd) {
		return nil
	}

	// Get arguments without flags, excluding "install" itself
	args := shell.GetArgsNoFlags(cmd)
	var packages []string
	for _, arg := range args {
		if arg != "install" {
			packages = append(packages, arg)
		}
	}

	return packages
}

func versionFixed(pkg string) bool {
	// Package is version-fixed if it has:
	// 1. = (version pin like package=1.2.3)
	// 2. / (path like /path/to/package.deb)
	// 3. .deb suffix (debian package file)
	return strings.Contains(pkg, "=") ||
		strings.Contains(pkg, "/") ||
		strings.HasSuffix(pkg, ".deb")
}
