package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3035 checks for zypper dist-upgrade usage.
func DL3035() rule.Rule {
	return rule.NewSimpleRule(
		DL3035Meta.Code,
		DL3035Meta.Severity,
		DL3035Meta.Message,
		checkDL3035,
	)
}

func checkDL3035(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Fail if any command uses zypper dist-upgrade or dup
	for _, cmd := range parsed.PresentCommands {
		if shell.CmdHasArgs("zypper", []string{"dist-upgrade"}, cmd) ||
			shell.CmdHasArgs("zypper", []string{"dup"}, cmd) {
			return false
		}
	}

	return true
}
