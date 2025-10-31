package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL4005 checks for ln /bin/sh instead of using SHELL.
func DL4005() rule.Rule {
	return rule.NewSimpleRule(
		DL4005Meta.Code,
		DL4005Meta.Severity,
		DL4005Meta.Message,
		checkDL4005,
	)
}

func checkDL4005(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Fail if using `ln /bin/sh` command
	for _, cmd := range parsed.PresentCommands {
		if shell.CmdHasArgs("ln", []string{"/bin/sh"}, cmd) {
			return false
		}
	}

	return true
}
