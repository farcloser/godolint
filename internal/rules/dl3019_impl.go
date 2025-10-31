package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3019 checks for apk add without --no-cache flag.
func DL3019() rule.Rule {
	return rule.NewSimpleRule(
		DL3019Meta.Code,
		DL3019Meta.Severity,
		DL3019Meta.Message,
		checkDL3019,
	)
}

func checkDL3019(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	// Skip if has cache/tmpfs mount for /var/cache/apk
	if hasCacheOrTmpfsMount(run.Flags, "/var/cache/apk") {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check if any command is apk add without --no-cache
	for _, cmd := range parsed.PresentCommands {
		if forgotApkNoCacheOption(cmd) {
			return false
		}
	}

	return true
}

func forgotApkNoCacheOption(cmd shell.Command) bool {
	// Must be apk add
	if !shell.CmdHasArgs("apk", []string{"add"}, cmd) {
		return false
	}

	// Check if it has --no-cache flag
	return !shell.HasFlag("no-cache", cmd)
}
