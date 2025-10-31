package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3046 checks for useradd without -l flag with high UID.
func DL3046() rule.Rule {
	return rule.NewSimpleRule(
		DL3046Meta.Code,
		DL3046Meta.Severity,
		DL3046Meta.Message,
		checkDL3046,
	)
}

func checkDL3046(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all useradd commands
	for _, cmd := range parsed.PresentCommands {
		if forgotUseraddFlagL(cmd) {
			return false
		}
	}

	return true
}

func forgotUseraddFlagL(cmd shell.Command) bool {
	// Must be useradd
	if cmd.Name != "useradd" {
		return false
	}

	// Must not have -l or --no-log-init flag
	hasLFlag := shell.HasAnyFlag([]string{"l", "no-log-init"}, cmd)
	if hasLFlag {
		return false
	}

	// Must have -u or --uid flag
	hasUFlag := shell.HasAnyFlag([]string{"u", "uid"}, cmd)
	if !hasUFlag {
		return false
	}

	// Check if UID is long (> 5 digits)
	uids := shell.GetFlagArg("u", cmd)
	for _, uid := range uids {
		if len(uid) > 5 {
			return true
		}
	}

	return false
}
