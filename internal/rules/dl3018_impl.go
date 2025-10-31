package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3018 checks for apk add without version pinning.
func DL3018() rule.Rule {
	return rule.NewSimpleRule(
		DL3018Meta.Code,
		DL3018Meta.Severity,
		DL3018Meta.Message,
		checkDL3018,
	)
}

func checkDL3018(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all apk add commands
	for _, cmd := range parsed.PresentCommands {
		packages := apkAddPackages(cmd)
		for _, pkg := range packages {
			if !isApkVersionFixed(pkg) {
				return false
			}
		}
	}

	return true
}

func apkAddPackages(cmd shell.Command) []string {
	// Only check apk add commands
	if !shell.CmdHasArgs("apk", []string{"add"}, cmd) {
		return nil
	}

	// Flags that take an argument we should skip
	flagsWithArgs := map[string]bool{
		"t":          true, // -t <target>
		"virtual":    true, // --virtual <name>
		"repository": true, // --repository <url>
		"X":          true, // -X <url>
	}

	// Get flag IDs to skip
	skipNextArgIDs := make(map[int]bool)

	for _, flag := range cmd.Flags {
		flagName := strings.TrimPrefix(flag.Arg, "-")
		flagName = strings.TrimPrefix(flagName, "-")

		// Check if this flag takes an argument
		if flagsWithArgs[flagName] {
			// Mark the next argument ID to skip
			// The next ID after this flag's ID is the argument value
			skipNextArgIDs[flag.ID+1] = true
		}
	}

	// Also skip flags themselves
	flagIDs := make(map[int]bool)
	for _, flag := range cmd.Flags {
		flagIDs[flag.ID] = true
	}

	// Get arguments, excluding flags and their values
	var packages []string

	for _, arg := range cmd.Arguments {
		// Skip if it's a flag
		if flagIDs[arg.ID] {
			continue
		}
		// Skip if it's a flag argument value
		if skipNextArgIDs[arg.ID] {
			continue
		}
		// Skip "add" itself
		if arg.Arg == "add" {
			continue
		}

		packages = append(packages, arg.Arg)
	}

	return packages
}

func isApkVersionFixed(pkg string) bool {
	// Package is version-fixed if it has:
	// 1. = (version pin like package=1.2.3)
	// 2. .apk suffix (package file)
	return strings.Contains(pkg, "=") || strings.HasSuffix(pkg, ".apk")
}
