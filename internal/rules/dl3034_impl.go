package rules

import (
	"slices"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// zypperCommand is the zypper package-manager binary, matched by the DL303x rules.
const zypperCommand = "zypper"

// DL3034 checks for zypper commands without non-interactive flag.
func DL3034() rule.Rule {
	return rule.NewSimpleRule(
		DL3034Meta.Code,
		DL3034Meta.Severity,
		DL3034Meta.Message,
		checkDL3034,
	)
}

func checkDL3034(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all zypper commands
	return !slices.ContainsFunc(parsed.PresentCommands, forgotZypperYesOption)
}

func forgotZypperYesOption(cmd shell.Command) bool {
	if !isZypperInstall(cmd) {
		return false
	}

	return !hasZypperYesOption(cmd)
}

func isZypperInstall(cmd shell.Command) bool {
	return shell.CmdHasArgs(zypperCommand, []string{"install"}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"in"}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"remove"}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"rm"}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"source-install"}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"si"}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"patch"}, cmd)
}

func hasZypperYesOption(cmd shell.Command) bool {
	return shell.HasAnyFlag([]string{"non-interactive", "n", "no-confirm", "y"}, cmd)
}
