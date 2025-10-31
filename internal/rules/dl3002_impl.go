package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3002State tracks root USER instructions per stage.
// Ported from Acc in DL3002.hs.
type dl3002State struct {
	currentStage int        // line number of current FROM
	rootUsers    map[int]int // map[stageLine]userLine - tracks last root USER per stage
}

// DL3002Rule checks that the last USER in each stage is not root.
// Ported from Hadolint.Rule.DL3002.
type DL3002Rule struct{}

// DL3002 creates the rule for checking last USER should not be root.
func DL3002() rule.Rule {
	return &DL3002Rule{}
}

func (r *DL3002Rule) Code() rule.RuleCode {
	return DL3002Meta.Code
}

func (r *DL3002Rule) Severity() rule.Severity {
	return DL3002Meta.Severity
}

func (r *DL3002Rule) Message() string {
	return DL3002Meta.Message
}

func (r *DL3002Rule) InitialState() rule.State {
	return rule.EmptyState(dl3002State{
		currentStage: -1,
		rootUsers:    make(map[int]int),
	})
}

// Check tracks USER instructions and remembers which stages end with root.
// Ported from the check function in DL3002.hs.
func (r *DL3002Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3002State)

	// Remember new stage
	if _, ok := instruction.(*syntax.From); ok {
		s.currentStage = line
		return state.ReplaceData(s)
	}

	// Track USER instruction
	if user, ok := instruction.(*syntax.User); ok {
		if isRoot(user.User) {
			// Root user - remember this line for current stage
			newRootUsers := make(map[int]int)
			for k, v := range s.rootUsers {
				newRootUsers[k] = v
			}
			newRootUsers[s.currentStage] = line
			s.rootUsers = newRootUsers
			return state.ReplaceData(s)
		}

		// Non-root user - forget stage (clear any previous root USER)
		newRootUsers := make(map[int]int)
		for k, v := range s.rootUsers {
			if k != s.currentStage {
				newRootUsers[k] = v
			}
		}
		s.rootUsers = newRootUsers
		return state.ReplaceData(s)
	}

	return state
}

// Finalize adds failures for all stages that end with a root USER.
// Ported from markFailures in DL3002.hs.
func (r *DL3002Rule) Finalize(state rule.State) rule.State {
	s := state.Data.(dl3002State)

	// Add failures for all stages with root users
	finalState := state
	for _, userLine := range s.rootUsers {
		finalState = finalState.AddFailure(rule.CheckFailure{
			Code:     DL3002Meta.Code,
			Severity: DL3002Meta.Severity,
			Message:  DL3002Meta.Message,
			Line:     userLine,
		})
	}

	return finalState
}

// isRoot checks if a USER instruction specifies root.
// Ported from isRoot in DL3002.hs.
func isRoot(user string) bool {
	return strings.HasPrefix(user, "root:") ||
		strings.HasPrefix(user, "0:") ||
		user == "root" ||
		user == "0"
}
