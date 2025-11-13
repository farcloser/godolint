package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3010State tracks copied archives and which ones get extracted.
type dl3010State struct {
	archives  map[string]int // map[basename]line - tracks copied archive files
	extracted map[string]int // map[basename]line - tracks which archives were extracted
}

// DL3010Rule checks for using COPY instead of ADD for archive extraction.
type DL3010Rule struct{}

// DL3010 creates the rule for checking ADD usage for archives.
func DL3010() rule.Rule {
	return &DL3010Rule{}
}

// Code returns the rule code.
func (*DL3010Rule) Code() rule.RuleCode {
	return DL3010Meta.Code
}

// Severity returns the rule severity.
func (*DL3010Rule) Severity() rule.Severity {
	return DL3010Meta.Severity
}

// Message returns the rule message.
func (*DL3010Rule) Message() string {
	return DL3010Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3010Rule) InitialState() rule.State {
	return rule.EmptyState(dl3010State{
		archives:  make(map[string]int),
		extracted: make(map[string]int),
	})
}

func (r *DL3010Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3010State)

	switch inst := instruction.(type) {
	case *syntax.From:
		// Reset state for new stage
		return rule.EmptyState(dl3010State{
			archives:  make(map[string]int),
			extracted: make(map[string]int),
		})

	case *syntax.Copy:
		// Track archive files being copied
		destBase := basename(inst.Destination)

		// If destination looks like an archive file, track its basename
		if isArchive(destBase) {
			s.archives[destBase] = line
		} else {
			// If destination is a directory, track source basenames
			for _, src := range inst.Source {
				srcBase := basename(src)
				if isArchive(srcBase) {
					s.archives[srcBase] = line
				}
			}
		}

		return state.ReplaceData(s)

	case *syntax.Run:
		parsed, err := shell.ParseShell(inst.Command)
		if err != nil {
			return state
		}

		// Check if any tracked archives are extracted
		for _, cmd := range parsed.PresentCommands {
			if isTarExtractCommand(cmd) || isUnzipCommand(cmd) {
				args := shell.GetArgsNoFlags(cmd)
				for _, arg := range args {
					base := basename(arg)
					if _, tracked := s.archives[base]; tracked {
						s.extracted[base] = s.archives[base] // Store original COPY line
					}
				}
			}
		}

		return state.ReplaceData(s)
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*DL3010Rule) Finalize(state rule.State) rule.State {
	s := state.Data.(dl3010State)

	finalState := state

	// Add failures for all extracted archives
	for _, line := range s.extracted {
		finalState = finalState.AddFailure(rule.CheckFailure{
			Code:     DL3010Meta.Code,
			Severity: DL3010Meta.Severity,
			Message:  DL3010Meta.Message,
			Line:     line,
		})
	}

	return finalState
}

func basename(path string) string {
	// Remove quotes
	path = dropQuotes(path)

	// Get last component after / or \
	idx := strings.LastIndexAny(path, "/\\")
	if idx != -1 {
		return path[idx+1:]
	}

	return path
}

// isArchive is defined in dl3020_impl.go and shared

func isTarExtractCommand(cmd shell.Command) bool {
	if cmd.Name != "tar" {
		return false
	}

	args := shell.GetArgs(cmd)

	// Check for long extract flags
	for _, arg := range args {
		if arg == "--extract" || arg == "--get" {
			return true
		}
	}

	// Check for short extract flags (-x)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && strings.Contains(arg, "x") {
			return true
		}
	}

	return false
}

var unzipCommands = map[string]bool{
	"unzip":      true,
	"gunzip":     true,
	"bunzip2":    true,
	"unlzma":     true,
	"unxz":       true,
	"zgz":        true,
	"uncompress": true,
	"zcat":       true,
	"gzcat":      true,
}

func isUnzipCommand(cmd shell.Command) bool {
	return unzipCommands[cmd.Name]
}
