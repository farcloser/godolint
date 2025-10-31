package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

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
	for _, cmd := range parsed.PresentCommands {
		if forgotZypperYesOption(cmd) {
			return false
		}
	}

	return true
}

func forgotZypperYesOption(cmd shell.Command) bool {
	if !isZypperInstall(cmd) {
		return false
	}

	return !hasZypperYesOption(cmd)
}

func isZypperInstall(cmd shell.Command) bool {
	return shell.CmdHasArgs("zypper", []string{"install"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"in"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"remove"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"rm"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"source-install"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"si"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"patch"}, cmd)
}

func hasZypperYesOption(cmd shell.Command) bool {
	return shell.HasAnyFlag([]string{"non-interactive", "n", "no-confirm", "y"}, cmd)
}
