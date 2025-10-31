package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3062 checks for go install/get/run without version pinning.
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

	// Check all go packages
	packages := getGoPackages(parsed)
	for _, pkg := range packages {
		if !isGoVersionPinned(pkg) {
			return false
		}
	}

	return true
}

var goCommands = []string{"install", "get", "run"}

func getGoPackages(parsed *shell.ParsedShell) []string {
	var packages []string

	for _, cmd := range parsed.PresentCommands {
		if !isGoCommand(cmd) {
			continue
		}

		// Get packages from arguments
		args := shell.GetArgsNoFlags(cmd)
		for _, arg := range args {
			// Skip command names
			if arg == "install" || arg == "get" || arg == "run" {
				continue
			}
			packages = append(packages, arg)
		}
	}

	return packages
}

func isGoCommand(cmd shell.Command) bool {
	if cmd.Name != "go" {
		return false
	}

	for _, goCmd := range goCommands {
		if shell.CmdHasArgs("go", []string{goCmd}, cmd) {
			return true
		}
	}

	return false
}

func isGoVersionPinned(pkg string) bool {
	// Must have @ symbol
	if !strings.Contains(pkg, "@") {
		return false
	}

	// Must not end with @latest or @none
	if strings.HasSuffix(pkg, "@latest") || strings.HasSuffix(pkg, "@none") {
		return false
	}

	return true
}
