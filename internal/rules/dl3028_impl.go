package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3028 checks for gem install without version pinning.
func DL3028() rule.Rule {
	return rule.NewSimpleRule(
		DL3028Meta.Code,
		DL3028Meta.Severity,
		DL3028Meta.Message,
		checkDL3028,
	)
}

func checkDL3028(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all gem install commands
	for _, cmd := range parsed.PresentCommands {
		gems := getGemPackages(cmd)
		for _, gem := range gems {
			if !isGemVersionFixed(gem) {
				return false
			}
		}
	}

	return true
}

func getGemPackages(cmd shell.Command) []string {
	// Must be gem install or gem i
	if !shell.CmdHasArgs("gem", []string{"install"}, cmd) && !shell.CmdHasArgs("gem", []string{"i"}, cmd) {
		return nil
	}

	// Skip if using -v or --version flag
	args := shell.GetArgs(cmd)
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			return nil
		}

		if strings.HasPrefix(arg, "--version=") {
			return nil
		}
	}

	// Get packages from arguments
	// Process only args until "--" separator
	argsUntilDoubleDash := []string{}

	for _, arg := range args {
		if arg == "--" {
			break
		}

		argsUntilDoubleDash = append(argsUntilDoubleDash, arg)
	}

	// Remove flags and their values
	var packages []string

	skipNext := false

	for _, arg := range argsUntilDoubleDash {
		if skipNext {
			skipNext = false

			continue
		}

		// Skip "install" and "i" commands
		if arg == "install" || arg == "i" {
			continue
		}

		// If it's a flag, skip it and next arg
		if strings.HasPrefix(arg, "-") {
			// For flags like --foo or -f, skip the next argument too
			if !strings.Contains(arg, "=") {
				skipNext = true
			}

			continue
		}

		packages = append(packages, arg)
	}

	return packages
}

func isGemVersionFixed(pkg string) bool {
	return strings.Contains(pkg, ":")
}
