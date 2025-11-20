package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl4006State tracks whether pipefail is set.
type dl4006State struct {
	pipefailSet bool
}

// DL4006Rule checks for pipefail with pipes.
type DL4006Rule struct{}

// DL4006 creates the rule for checking pipefail.
func DL4006() rule.Rule {
	return &DL4006Rule{}
}

// Code returns the rule code.
func (*DL4006Rule) Code() rule.RuleCode {
	return DL4006Meta.Code
}

// Severity returns the rule severity.
func (*DL4006Rule) Severity() rule.Severity {
	return DL4006Meta.Severity
}

// Message returns the rule message.
func (*DL4006Rule) Message() string {
	return DL4006Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL4006Rule) InitialState() rule.State {
	return rule.EmptyState(dl4006State{
		pipefailSet: false,
	})
}

func (r *DL4006Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl4006State)

	switch inst := instruction.(type) {
	case *syntax.From:
		// Reset state on new FROM
		s.pipefailSet = false

		return state.ReplaceData(s)

	case *syntax.Shell:
		// Check if SHELL sets pipefail
		if len(inst.Arguments) > 0 {
			shellCmd := strings.Join(inst.Arguments, " ")

			// Check if it's a non-POSIX shell (fish, powershell, etc)
			if isNonPosixShell(shellCmd) {
				s.pipefailSet = true // Skip checks for non-POSIX shells

				return state.ReplaceData(s)
			}

			// Check if pipefail is set
			s.pipefailSet = hasPipefailOption(shellCmd)
		}

		return state.ReplaceData(s)

	case *syntax.Run:
		// If pipefail is not set and command has pipes, fail
		if !s.pipefailSet && hasPipes(inst.Command) {
			return state.AddFailure(rule.CheckFailure{
				Code:     DL4006Meta.Code,
				Severity: DL4006Meta.Severity,
				Message:  DL4006Meta.Message,
				Line:     line,
				Column:   1, // Hardcoded to 1 (matches hadolint)
			})
		}

		return state
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*DL4006Rule) Finalize(state rule.State) rule.State {
	return state
}

// isNonPosixShell checks if shell is non-POSIX (fish, powershell, etc).
func isNonPosixShell(shellCmd string) bool {
	nonPosixShells := []string{
		"/usr/bin/fish",
		"/bin/fish",
		"fish",
		"/usr/bin/pwsh",
		"/bin/pwsh",
		"pwsh",
		"powershell",
		"cmd.exe",
	}

	for _, shell := range nonPosixShells {
		if strings.Contains(shellCmd, shell) {
			return true
		}
	}

	return false
}

// hasPipes checks if a command contains pipes.
func hasPipes(command string) bool {
	// Parse the shell script to check for actual pipe operators
	parsed, err := shell.ParseShell(command)
	if err != nil {
		// Can't parse = can't reliably detect pipes
		return false
	}

	// Use shell package to detect pipes properly
	return shell.HasPipes(parsed)
}

// hasPipefailOption checks if a shell command sets pipefail.
func hasPipefailOption(shellCmd string) bool {
	// Parse the shell command to check for -o pipefail
	parsed, err := shell.ParseShell(shellCmd)
	if err != nil {
		return false
	}

	validShells := []string{
		"/bin/bash",
		"/bin/zsh",
		"/bin/ash",
		"bash",
		"zsh",
		"ash",
	}

	for _, cmd := range parsed.PresentCommands {
		// Check if it's a valid shell
		isValidShell := false

		for _, validShell := range validShells {
			if cmd.Name == validShell {
				isValidShell = true

				break
			}
		}

		if !isValidShell {
			continue
		}

		// Check for -o flag
		if !shell.HasFlag("o", cmd) {
			continue
		}

		// Check if pipefail is in arguments
		args := shell.GetArgs(cmd)
		for _, arg := range args {
			if arg == "pipefail" {
				return true
			}
		}
	}

	return false
}
