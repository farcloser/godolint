package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3040 checks for dnf clean all after dnf install.
func DL3040() rule.Rule {
	return rule.NewSimpleRule(
		DL3040Meta.Code,
		DL3040Meta.Severity,
		DL3040Meta.Message,
		checkDL3040,
	)
}

func checkDL3040(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	// Check if cache/tmpfs mount is present
	if hasCacheOrTmpfsMount(run.Flags, "/var/cache/libdnf5") ||
		hasCacheOrTmpfsMount(run.Flags, ".cache/libdnf5") {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	hasDnfInstall := false
	hasMicroDnfInstall := false
	hasDnfClean := false
	hasMicroDnfClean := false

	for _, cmd := range parsed.PresentCommands {
		if shell.CmdHasArgs("dnf", []string{"install"}, cmd) {
			hasDnfInstall = true
		}
		if shell.CmdHasArgs("microdnf", []string{"install"}, cmd) {
			hasMicroDnfInstall = true
		}
		if isDnfCleanCmd(cmd) {
			hasDnfClean = true
		}
		if isMicroDnfCleanCmd(cmd) {
			hasMicroDnfClean = true
		}
	}

	// If has dnf install, must have dnf clean
	if hasDnfInstall && !hasDnfClean {
		return false
	}

	// If has microdnf install, must have microdnf clean
	if hasMicroDnfInstall && !hasMicroDnfClean {
		return false
	}

	return true
}

func isDnfCleanCmd(cmd shell.Command) bool {
	if shell.CmdHasArgs("dnf", []string{"clean", "all"}, cmd) {
		return true
	}
	if shell.CmdHasArgs("rm", []string{"-rf", "/var/cache/libdnf5*"}, cmd) {
		return true
	}
	return false
}

func isMicroDnfCleanCmd(cmd shell.Command) bool {
	if shell.CmdHasArgs("microdnf", []string{"clean", "all"}, cmd) {
		return true
	}
	if shell.CmdHasArgs("rm", []string{"-rf", "/var/cache/libdnf5*"}, cmd) {
		return true
	}
	return false
}
