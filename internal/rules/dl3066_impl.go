package rules

import (
	"maps"
	"strings"
	"unicode"

	"github.com/forkcloser/godolint/internal/rule"
	"github.com/forkcloser/godolint/internal/syntax"
)

// dl3066State is the set of ARG names declared so far: a user id written as
// one of them is resolved at build time, not a name the host has to know.
type dl3066State struct {
	args map[string]bool
}

// DL3066Rule checks that a USER names a numeric id.
// Ported from Hadolint.Rule.DL3066.
type DL3066Rule struct{}

// DL3066 creates the rule that a user id is numeric, or a declared ARG.
func DL3066() rule.Rule {
	return &DL3066Rule{}
}

// Code returns the rule code.
func (*DL3066Rule) Code() rule.Code {
	return DL3066Meta.Code
}

// Severity returns the rule severity.
func (*DL3066Rule) Severity() rule.Severity {
	return DL3066Meta.Severity
}

// Message returns the rule message.
func (*DL3066Rule) Message() string {
	return DL3066Meta.Message
}

// InitialState returns the initial state for this rule.
func (*DL3066Rule) InitialState() rule.State {
	return rule.EmptyState(dl3066State{args: make(map[string]bool)})
}

// Check registers each ARG, and fails a USER whose id (the part before any
// `:`) is neither all digits nor a reference to a registered ARG.
func (*DL3066Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	switch instr := instruction.(type) {
	case *syntax.Arg:
		current := rule.Data[dl3066State](state)
		args := make(map[string]bool, len(current.args)+1)
		maps.Copy(args, current.args)
		args[instr.ArgName] = true

		return state.ReplaceData(dl3066State{args: args})
	case *syntax.User:
		uid, _, _ := strings.Cut(instr.User, ":")
		if allDigits(uid) || uidIsDefinedArg(rule.Data[dl3066State](state), uid) {
			return state
		}

		return state.AddFailure(rule.CheckFailure{
			Code:     DL3066Meta.Code,
			Severity: DL3066Meta.Severity,
			Message:  DL3066Meta.Message,
			Line:     line,
			Column:   1,
		})
	default:
		return state
	}
}

// Finalize performs no final check.
func (*DL3066Rule) Finalize(state rule.State) rule.State {
	return state
}

// allDigits is hadolint's `Text.all Char.isDigit`: true of the empty string.
func allDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// uidIsDefinedArg reports whether uid refers to a registered ARG, as `${NAME}`
// or `$NAME`.
func uidIsDefinedArg(current dl3066State, uid string) bool {
	for arg := range current.args {
		if strings.Contains(uid, "${"+arg+"}") || strings.Contains(uid, "$"+arg) {
			return true
		}
	}

	return false
}
