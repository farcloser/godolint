package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3015 checks for apt-get install without --no-install-recommends.
func DL3015() rule.Rule {
	return rule.NewSimpleRule(
		DL3015Meta.Code,
		DL3015Meta.Severity,
		DL3015Meta.Message,
		checkDL3015,
	)
}

func checkDL3015(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check if any command is apt-get install without --no-install-recommends
	for _, cmd := range parsed.PresentCommands {
		if forgotNoInstallRecommends(cmd) {
			return false
		}
	}

	return true
}

func forgotNoInstallRecommends(cmd shell.Command) bool {
	// Must be apt-get install
	if !shell.CmdHasArgs("apt-get", []string{"install"}, cmd) {
		return false
	}

	// Check if it disables recommend option
	return !disablesRecommendOption(cmd)
}

func disablesRecommendOption(cmd shell.Command) bool {
	// Check for --no-install-recommends flag
	if shell.HasFlag("no-install-recommends", cmd) {
		return true
	}

	// Check for APT::Install-Recommends=false argument
	if shell.HasArg("APT::Install-Recommends=false", cmd) {
		return true
	}

	return false
}
