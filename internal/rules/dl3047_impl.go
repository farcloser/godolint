package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3047 checks for wget without progress bar option.
func DL3047() rule.Rule {
	return rule.NewSimpleRule(
		DL3047Meta.Code,
		DL3047Meta.Severity,
		DL3047Meta.Message,
		checkDL3047,
	)
}

func checkDL3047(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all wget commands
	for _, cmd := range parsed.PresentCommands {
		if forgotWgetProgress(cmd) {
			return false
		}
	}

	return true
}

func forgotWgetProgress(cmd shell.Command) bool {
	// Must be wget
	if cmd.Name != "wget" {
		return false
	}

	// If has progress option, didn't forget
	if shell.HasFlag("progress", cmd) {
		return false
	}

	// If has special flags that suppress output, didn't forget
	if hasWgetSpecialFlags(cmd) {
		return false
	}

	// Missing progress option and no special flags - forgot!
	return true
}

func hasWgetSpecialFlags(cmd shell.Command) bool {
	// Quiet flags
	if shell.HasAnyFlag([]string{"q", "quiet"}, cmd) {
		return true
	}

	// Output redirection flags
	if shell.HasAnyFlag([]string{"o", "output-file"}, cmd) {
		return true
	}

	// Append output flags
	if shell.HasAnyFlag([]string{"a", "append-output"}, cmd) {
		return true
	}

	// No-verbose flag
	if shell.HasFlag("no-verbose", cmd) {
		return true
	}

	// Special case: -nv (combined flag)
	if shell.CmdHasArgs("wget", []string{"-nv"}, cmd) {
		return true
	}

	return false
}
