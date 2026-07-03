package rules

import (
	"slices"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// yumCommand is the yum package-manager binary, matched by the DL303x rules.
const yumCommand = "yum"

// DL3030 checks for yum install without -y flag.
func DL3030() rule.Rule {
	return rule.NewSimpleRule(
		DL3030Meta.Code,
		DL3030Meta.Severity,
		DL3030Meta.Message,
		checkDL3030,
	)
}

func checkDL3030(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all yum install commands
	return !slices.ContainsFunc(parsed.PresentCommands, forgotYumYesOption)
}

func forgotYumYesOption(cmd shell.Command) bool {
	// Must be yum install/groupinstall/localinstall
	if !isYumInstall(cmd) {
		return false
	}

	// Must have -y or --assumeyes flag
	return !hasYumYesOption(cmd)
}

func isYumInstall(cmd shell.Command) bool {
	return shell.CmdHasArgs(yumCommand, []string{"install"}, cmd) ||
		shell.CmdHasArgs(yumCommand, []string{"groupinstall"}, cmd) ||
		shell.CmdHasArgs(yumCommand, []string{"localinstall"}, cmd)
}

func hasYumYesOption(cmd shell.Command) bool {
	return shell.HasAnyFlag([]string{"y", "assumeyes"}, cmd)
}
