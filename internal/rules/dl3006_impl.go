package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3006State tracks FROM aliases to allow untagged images that reference aliases.
type dl3006State struct {
	aliases map[string]bool // Set of defined FROM aliases
}

// DL3006Rule checks that images have explicit tags.
// Ported from Hadolint.Rule.DL3006.
type DL3006Rule struct{}

// DL3006 creates the rule for checking images have explicit tags.
func DL3006() rule.Rule {
	return &DL3006Rule{}
}

func (r *DL3006Rule) Code() rule.RuleCode {
	return DL3006Meta.Code
}

func (r *DL3006Rule) Severity() rule.Severity {
	return DL3006Meta.Severity
}

func (r *DL3006Rule) Message() string {
	return DL3006Meta.Message
}

func (r *DL3006Rule) InitialState() rule.State {
	return rule.EmptyState(dl3006State{
		aliases: make(map[string]bool),
	})
}

// Check examines FROM instructions and tracks aliases.
// Ported from the check function in DL3006.hs.
func (r *DL3006Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	from, ok := instruction.(*syntax.From)
	if !ok {
		return state
	}

	s := state.Data.(dl3006State)

	// Add alias to set if present
	if from.Image.Alias != nil {
		newAliases := make(map[string]bool)
		for k, v := range s.aliases {
			newAliases[k] = v
		}

		newAliases[*from.Image.Alias] = true
		s.aliases = newAliases
		state = state.ReplaceData(s)
	}

	// Check if image needs explicit tag
	// Scratch image - OK
	if from.Image.Image == "scratch" {
		return state
	}

	// Has digest - OK
	if from.Image.Digest != nil {
		return state
	}

	// Has tag - OK
	if from.Image.Tag != nil {
		return state
	}

	// Variable reference - OK
	if strings.HasPrefix(from.Image.Image, "$") {
		return state
	}

	// FROM alias reference - OK
	if s.aliases[from.Image.Image] {
		return state
	}

	// No tag and not an exception - fail
	return state.AddFailure(rule.CheckFailure{
		Code:     DL3006Meta.Code,
		Severity: DL3006Meta.Severity,
		Message:  DL3006Meta.Message,
		Line:     line,
	})
}

func (r *DL3006Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
