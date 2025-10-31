package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3009State tracks apt list cleanup per stage.
type dl3009State struct {
	lastFrom    *syntax.BaseImage // Last FROM instruction
	dockerClean bool              // Whether docker-clean is enabled
	stages      map[string]int    // map[alias]line - tracks stage names
	forgets     map[int]string    // map[line]alias - tracks lines that forgot cleanup
}

// DL3009Rule checks for deletion of apt lists after apt update.
type DL3009Rule struct{}

// DL3009 creates the rule for checking apt list cleanup.
func DL3009() rule.Rule {
	return &DL3009Rule{}
}

func (r *DL3009Rule) Code() rule.RuleCode {
	return DL3009Meta.Code
}

func (r *DL3009Rule) Severity() rule.Severity {
	return DL3009Meta.Severity
}

func (r *DL3009Rule) Message() string {
	return DL3009Meta.Message
}

func (r *DL3009Rule) InitialState() rule.State {
	return rule.EmptyState(dl3009State{
		dockerClean: true,
		stages:      make(map[string]int),
		forgets:     make(map[int]string),
	})
}

func (r *DL3009Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3009State)

	switch inst := instruction.(type) {
	case *syntax.From:
		// Remember new stage
		s.lastFrom = &inst.Image
		s.dockerClean = true
		// Track which image names are referenced by FROM instructions
		// This is used to detect if a stage (by its alias) is used later
		s.stages[inst.Image.Image] = line

		return state.ReplaceData(s)

	case *syntax.Run:
		parsed, err := shell.ParseShell(inst.Command)
		if err != nil {
			return state
		}

		// Check if forgot to cleanup apt lists
		if forgotToCleanup(parsed) {
			// Skip if has cache/tmpfs mount for /var/lib/apt/lists
			if hasCacheOrTmpfsMount(inst.Flags, "/var/lib/apt/lists") {
				return state
			}

			// Skip if has cache/tmpfs mount for BOTH /var/lib/apt AND /var/cache/apt
			if hasCacheOrTmpfsMount(inst.Flags, "/var/lib/apt") &&
				hasCacheOrTmpfsMount(inst.Flags, "/var/cache/apt") {
				return state
			}

			// Record this line as forgetting cleanup
			alias := ""
			if s.lastFrom != nil && s.lastFrom.Alias != nil {
				alias = *s.lastFrom.Alias
			}

			s.forgets[line] = alias
		} else if disabledDockerClean(parsed) {
			s.dockerClean = false
		}

		return state.ReplaceData(s)
	}

	return state
}

func (r *DL3009Rule) Finalize(state rule.State) rule.State {
	s := state.Data.(dl3009State)

	finalState := state

	lastAlias := ""
	if s.lastFrom != nil && s.lastFrom.Alias != nil {
		lastAlias = *s.lastFrom.Alias
	}

	// Add failures for forgets that matter
	for line, alias := range s.forgets {
		// Fail if this is the last stage
		if alias == lastAlias {
			finalState = finalState.AddFailure(rule.CheckFailure{
				Code:     DL3009Meta.Code,
				Severity: DL3009Meta.Severity,
				Message:  DL3009Meta.Message,
				Line:     line,
			})

			continue
		}

		// Fail if this stage is used later (alias appears in stages)
		if alias != "" {
			if _, used := s.stages[alias]; used {
				finalState = finalState.AddFailure(rule.CheckFailure{
					Code:     DL3009Meta.Code,
					Severity: DL3009Meta.Severity,
					Message:  DL3009Meta.Message,
					Line:     line,
				})
			}
		}
	}

	return finalState
}

func forgotToCleanup(parsed *shell.ParsedShell) bool {
	hasUpdate := false
	hasCleanup := false

	for _, cmd := range parsed.PresentCommands {
		// Check for apt/apt-get/aptitude update
		if shell.CmdHasArgs("apt", []string{"update"}, cmd) ||
			shell.CmdHasArgs("apt-get", []string{"update"}, cmd) ||
			shell.CmdHasArgs("aptitude", []string{"update"}, cmd) {
			hasUpdate = true
		}

		// Check for cleanup
		if shell.CmdHasArgs("rm", []string{"-rf", "/var/lib/apt/lists/*"}, cmd) {
			hasCleanup = true
		}
	}

	return hasUpdate && !hasCleanup
}

func disabledDockerClean(parsed *shell.ParsedShell) bool {
	for _, cmd := range parsed.PresentCommands {
		// Check if removes docker-clean script
		if shell.CmdHasArgs("rm", []string{"/etc/apt/apt.conf.d/docker-clean"}, cmd) {
			return true
		}

		// Check if configures to keep packages
		if shell.CmdHasArgs("echo", []string{"'Binary::apt::APT::Keep-Downloaded-Packages \"true\";'"}, cmd) {
			return true
		}
	}

	return false
}
