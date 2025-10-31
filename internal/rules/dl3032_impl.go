package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3032 checks for yum clean all after yum install.
func DL3032() rule.Rule {
	return rule.NewSimpleRule(
		DL3032Meta.Code,
		DL3032Meta.Severity,
		DL3032Meta.Message,
		checkDL3032,
	)
}

func checkDL3032(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	hasYumInstall := false
	hasYumClean := false

	for _, cmd := range parsed.PresentCommands {
		if shell.CmdHasArgs("yum", []string{"install"}, cmd) {
			hasYumInstall = true
		}
		if isYumClean(cmd) {
			hasYumClean = true
		}
	}

	// If no yum install, pass
	if !hasYumInstall {
		return true
	}

	// If has yum install, must have clean
	return hasYumClean
}

func isYumClean(cmd shell.Command) bool {
	// yum clean all
	if shell.CmdHasArgs("yum", []string{"clean", "all"}, cmd) {
		return true
	}

	// rm -rf /var/cache/yum/*
	if shell.CmdHasArgs("rm", []string{"-rf", "/var/cache/yum/*"}, cmd) {
		return true
	}

	return false
}
