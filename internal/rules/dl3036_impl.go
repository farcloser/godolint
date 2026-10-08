package rules

import (
	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/shell"
	"github.com/forkcloser/godolint/internal/syntax"
)

// DL3036 checks for zypper clean after zypper install.
func DL3036() rule.Rule {
	return rule.NewSimpleRule(
		DL3036Meta.Code,
		DL3036Meta.Severity,
		DL3036Meta.Message,
		checkDL3036,
	)
}

// checkDL3036 passes a RUN with no zypper install, one whose cache is a
// mount, or one whose first zypper install comes before its first clean.
func checkDL3036(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	if firstCommand(parsed.PresentCommands, isZypperInstallCmd) < 0 {
		return true
	}

	if hasCacheOrTmpfsMount(run.Flags, "/var/cache/zypp") {
		return true
	}

	return cleanFollowsInstall(parsed.PresentCommands, isZypperInstallCmd, isZypperCleanCmd)
}

func isZypperInstallCmd(cmd shell.Command) bool {
	return shell.CmdHasArgs(zypperCommand, []string{installArg}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"in"}, cmd)
}

func isZypperCleanCmd(cmd shell.Command) bool {
	return shell.CmdHasArgs(zypperCommand, []string{cleanArg}, cmd) ||
		shell.CmdHasArgs(zypperCommand, []string{"cc"}, cmd)
}
