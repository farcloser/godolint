package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3024State tracks seen FROM aliases to detect duplicates.
type dl3024State struct {
	aliases map[string]int // map[alias]line - tracks where each alias was defined
}

// DL3024Rule checks that FROM aliases are unique.
// Ported from Hadolint.Rule.DL3024.
type DL3024Rule struct{}

// DL3024 creates the rule for checking FROM aliases are unique.
func DL3024() rule.Rule {
	return &DL3024Rule{}
}

// Code returns the rule code.
func (*DL3024Rule) Code() rule.RuleCode {
	return DL3024Meta.Code
}

// Severity returns the rule severity.
func (*DL3024Rule) Severity() rule.Severity {
	return DL3024Meta.Severity
}

// Message returns the rule message.
func (*DL3024Rule) Message() string {
	return DL3024Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3024Rule) InitialState() rule.State {
	return rule.EmptyState(dl3024State{
		aliases: make(map[string]int),
	})
}

// Check tracks FROM aliases and reports duplicates.
// Ported from the check function in DL3024.hs.
func (*DL3024Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	from, ok := instruction.(*syntax.From)
	if !ok {
		return state
	}

	// No alias - OK
	if from.Image.Alias == nil {
		return state
	}

	s := state.Data.(dl3024State)
	alias := *from.Image.Alias

	// Check if alias already seen
	if _, exists := s.aliases[alias]; exists {
		// Duplicate alias - fail
		return state.AddFailure(rule.CheckFailure{
			Code:     DL3024Meta.Code,
			Severity: DL3024Meta.Severity,
			Message:  DL3024Meta.Message,
			Line:     line,
			Column:   1, // Hardcoded to 1 (matches hadolint)
		})
	}

	// Remember this alias
	newAliases := make(map[string]int)
	for k, v := range s.aliases {
		newAliases[k] = v
	}

	newAliases[alias] = line
	s.aliases = newAliases

	return state.ReplaceData(s)
}

// Finalize performs final checks after processing all instructions.
func (*DL3024Rule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}
