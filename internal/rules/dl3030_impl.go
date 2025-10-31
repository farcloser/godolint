package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

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
	for _, cmd := range parsed.PresentCommands {
		if forgotYumYesOption(cmd) {
			return false
		}
	}

	return true
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
	return shell.CmdHasArgs("yum", []string{"install"}, cmd) ||
		shell.CmdHasArgs("yum", []string{"groupinstall"}, cmd) ||
		shell.CmdHasArgs("yum", []string{"localinstall"}, cmd)
}

func hasYumYesOption(cmd shell.Command) bool {
	return shell.HasAnyFlag([]string{"y", "assumeyes"}, cmd)
}
