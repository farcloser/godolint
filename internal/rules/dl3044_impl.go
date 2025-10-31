package rules

import (
	"regexp"
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3044State tracks defined environment variables.
type dl3044State struct {
	definedVars map[string]bool
}

// DL3044Rule checks for ENV self-reference.
type DL3044Rule struct{}

// DL3044 creates the rule for checking ENV self-reference.
func DL3044() rule.Rule {
	return &DL3044Rule{}
}

func (r *DL3044Rule) Code() rule.RuleCode {
	return DL3044Meta.Code
}

func (r *DL3044Rule) Severity() rule.Severity {
	return DL3044Meta.Severity
}

func (r *DL3044Rule) Message() string {
	return DL3044Meta.Message
}

func (r *DL3044Rule) InitialState() rule.State {
	return rule.EmptyState(dl3044State{
		definedVars: make(map[string]bool),
	})
}

func (r *DL3044Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3044State)

	switch inst := instruction.(type) {
	case *syntax.Arg:
		// Track ARG variables
		s.definedVars[inst.ArgName] = true

		return state.ReplaceData(s)

	case *syntax.Env:
		// Check if any variable in this ENV references another variable
		// defined in the same ENV statement
		newVars := make([]string, 0, len(inst.Pairs))
		for _, pair := range inst.Pairs {
			newVars = append(newVars, pair.Key)
		}

		// Check for self-references
		for i, pair := range inst.Pairs {
			// Check if this value references any of the other variables
			// defined in the same ENV (but not itself)
			for j, otherPair := range inst.Pairs {
				if i == j {
					continue // Skip self
				}

				// Check if value references the other variable
				if referencesVar(pair.Value, otherPair.Key) {
					// Only fail if the referenced variable is NOT already defined
					if !s.definedVars[otherPair.Key] {
						return state.AddFailure(rule.CheckFailure{
							Code:     DL3044Meta.Code,
							Severity: DL3044Meta.Severity,
							Message:  DL3044Meta.Message,
							Line:     line,
						})
					}
				}
			}
		}

		// Add all new variables to defined set
		for _, varName := range newVars {
			s.definedVars[varName] = true
		}

		return state.ReplaceData(s)
	}

	return state
}

func (r *DL3044Rule) Finalize(state rule.State) rule.State {
	return state
}

// referencesVar checks if a value string references a variable.
// Matches ${var} or $var (where var is terminated by non-alphanumeric char).
func referencesVar(value, varName string) bool {
	// Check for ${varName}
	if strings.Contains(value, "${"+varName+"}") {
		return true
	}

	// Check for $varName with termination
	// Match $varName where it's followed by non-variable character
	pattern := regexp.MustCompile(`\$` + regexp.QuoteMeta(varName) + `(?:[^a-zA-Z0-9_]|$)`)

	return pattern.MatchString(value)
}
