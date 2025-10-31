package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3001 checks for commands that make no sense in Docker containers.
func DL3001() rule.Rule {
	return rule.NewSimpleRule(
		DL3001Meta.Code,
		DL3001Meta.Severity,
		DL3001Meta.Message,
		checkDL3001,
	)
}

var invalidCommands = map[string]bool{
	"free":     true,
	"kill":     true,
	"mount":    true,
	"ps":       true,
	"service":  true,
	"shutdown": true,
	"ssh":      true,
	"top":      true,
	"vim":      true,
}

func checkDL3001(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	commands := shell.FindCommandNames(parsed)
	for _, cmd := range commands {
		if invalidCommands[cmd] {
			return false // Fail if invalid command found
		}
	}

	return true
}
