package rules

import (
	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/shell"
	"github.com/forkcloser/godolint/internal/syntax"
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

	var usage dnfUsage
	for _, cmd := range parsed.PresentCommands {
		usage.observe(cmd)
	}

	return usage.cleansWhatItInstalls()
}

// dnfUsage is what one RUN does with dnf and with microdnf. The two are
// tracked apart: each cleans its own cache, and cleaning the other's leaves
// the layer just as fat.
type dnfUsage struct {
	dnfInstall      bool
	dnfClean        bool
	microdnfInstall bool
	microdnfClean   bool
}

// observe folds one command of the RUN into the usage.
func (u *dnfUsage) observe(cmd shell.Command) {
	u.dnfInstall = u.dnfInstall || shell.CmdHasArgs(dnfCommand, []string{installArg}, cmd)
	u.microdnfInstall = u.microdnfInstall || shell.CmdHasArgs(microdnfCommand, []string{installArg}, cmd)
	u.dnfClean = u.dnfClean || isDnfCleanCmd(cmd)
	u.microdnfClean = u.microdnfClean || isMicroDnfCleanCmd(cmd)
}

// cleansWhatItInstalls reports whether every install in the RUN is followed by
// the matching clean.
func (u *dnfUsage) cleansWhatItInstalls() bool {
	return (!u.dnfInstall || u.dnfClean) && (!u.microdnfInstall || u.microdnfClean)
}

func isDnfCleanCmd(cmd shell.Command) bool {
	if shell.CmdHasArgs(dnfCommand, []string{cleanArg, allArg}, cmd) {
		return true
	}

	if shell.CmdHasArgs("rm", []string{recursiveForceFlag, "/var/cache/libdnf5*"}, cmd) {
		return true
	}

	return false
}

func isMicroDnfCleanCmd(cmd shell.Command) bool {
	if shell.CmdHasArgs(microdnfCommand, []string{cleanArg, allArg}, cmd) {
		return true
	}

	if shell.CmdHasArgs("rm", []string{recursiveForceFlag, "/var/cache/libdnf5*"}, cmd) {
		return true
	}

	return false
}
