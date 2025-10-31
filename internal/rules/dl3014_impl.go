package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3014 checks for apt-get install without -y flag.
func DL3014() rule.Rule {
	return rule.NewSimpleRule(
		DL3014Meta.Code,
		DL3014Meta.Severity,
		DL3014Meta.Message,
		checkDL3014,
	)
}

func checkDL3014(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check if any command is apt-get install without -y flag
	for _, cmd := range parsed.PresentCommands {
		if forgotAptYesOption(cmd) {
			return false
		}
	}

	return true
}

func forgotAptYesOption(cmd shell.Command) bool {
	// Must be apt-get install
	if !shell.CmdHasArgs("apt-get", []string{"install"}, cmd) {
		return false
	}

	// Check if it has yes option
	return !hasYesOption(cmd)
}

func hasYesOption(cmd shell.Command) bool {
	// Check for flags: -y, --yes, -qq, --assume-yes
	if shell.HasAnyFlag([]string{"y", "yes", "qq", "assume-yes"}, cmd) {
		return true
	}

	// Check for -q -q or --quiet --quiet
	if shell.CountFlag("q", cmd) == 2 || shell.CountFlag("quiet", cmd) == 2 {
		return true
	}

	// Check for -q=2
	if shell.HasArg("-q=2", cmd) {
		return true
	}

	return false
}
