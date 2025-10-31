package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3004 checks for use of `sudo`.
func DL3004() rule.Rule {
	return rule.NewSimpleRule(
		DL3004Meta.Code,
		DL3004Meta.Severity,
		DL3004Meta.Message,
		checkDL3004,
	)
}

func checkDL3004(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Fail if using `sudo` command
	return !shell.UsingProgram("sudo", parsed)
}
