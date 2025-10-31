package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
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

func checkDL3036(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	// Skip if has cache/tmpfs mount for /var/cache/zypp
	if hasCacheOrTmpfsMount(run.Flags, "/var/cache/zypp") {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	hasZypperInstall := false
	hasZypperClean := false

	for _, cmd := range parsed.PresentCommands {
		if isZypperInstallCmd(cmd) {
			hasZypperInstall = true
		}

		if isZypperCleanCmd(cmd) {
			hasZypperClean = true
		}
	}

	// If no zypper install, pass
	if !hasZypperInstall {
		return true
	}

	// If has zypper install, must have clean
	return hasZypperClean
}

func isZypperInstallCmd(cmd shell.Command) bool {
	return shell.CmdHasArgs("zypper", []string{"install"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"in"}, cmd)
}

func isZypperCleanCmd(cmd shell.Command) bool {
	return shell.CmdHasArgs("zypper", []string{"clean"}, cmd) ||
		shell.CmdHasArgs("zypper", []string{"cc"}, cmd)
}
