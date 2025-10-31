package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3057State tracks stages and their HEALTHCHECK status.
type dl3057State struct {
	currentStage *stageID
	goodStages   map[stageID]bool // stages with HEALTHCHECK or inherited
	badStages    map[stageID]bool // stages without HEALTHCHECK
}

type stageID struct {
	src  string // source image name (what this stage is FROM)
	name string // alias name (or same as src if no alias)
	line int    // line number where stage is defined
}

// DL3057Rule checks for missing HEALTHCHECK instructions.
type DL3057Rule struct{}

// DL3057 creates the rule for checking missing HEALTHCHECK.
func DL3057() rule.Rule {
	return &DL3057Rule{}
}

func (r *DL3057Rule) Code() rule.RuleCode {
	return DL3057Meta.Code
}

func (r *DL3057Rule) Severity() rule.Severity {
	return DL3057Meta.Severity
}

func (r *DL3057Rule) Message() string {
	return DL3057Meta.Message
}

func (r *DL3057Rule) InitialState() rule.State {
	return rule.EmptyState(dl3057State{
		currentStage: nil,
		goodStages:   make(map[stageID]bool),
		badStages:    make(map[stageID]bool),
	})
}

func (r *DL3057Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3057State)

	switch inst := instruction.(type) {
	case *syntax.From:
		// Create stage ID
		imageName := inst.Image.Image
		stageName := imageName
		if inst.Image.Alias != nil {
			stageName = *inst.Image.Alias
		}

		newStage := stageID{
			src:  imageName,
			name: stageName,
			line: line,
		}

		// Check if this stage inherits from a good stage
		inherited := false
		for goodStage := range s.goodStages {
			if goodStage.name == imageName {
				inherited = true
				break
			}
		}

		if inherited {
			// Mark as good since it inherits from a good stage
			s.goodStages[newStage] = true
		} else {
			// Mark as bad for now (can be updated if HEALTHCHECK found)
			s.badStages[newStage] = true
		}

		s.currentStage = &newStage
		return state.ReplaceData(s)

	case *syntax.Healthcheck:
		// Mark current stage and all its ancestors as good
		if s.currentStage != nil {
			s = markGood(s, *s.currentStage)
		}
		return state.ReplaceData(s)
	}

	return state
}

// markGood marks a stage and all its ancestors as good.
func markGood(s dl3057State, stage stageID) dl3057State {
	// Mark this stage as good
	s.goodStages[stage] = true
	delete(s.badStages, stage)

	// Find and mark ancestors recursively
	ancestors := findAncestors(s, stage)
	for ancestor := range ancestors {
		s.goodStages[ancestor] = true
		delete(s.badStages, ancestor)
	}

	return s
}

// findAncestors finds all ancestor stages recursively.
func findAncestors(s dl3057State, stage stageID) map[stageID]bool {
	ancestors := make(map[stageID]bool)

	// Find stages in badStages that this stage inherits from
	for badStage := range s.badStages {
		if badStage.name == stage.src {
			ancestors[badStage] = true
			// Recursively find ancestors of this ancestor
			for ancestorOfAncestor := range findAncestors(s, badStage) {
				ancestors[ancestorOfAncestor] = true
			}
		}
	}

	return ancestors
}

func (r *DL3057Rule) Finalize(state rule.State) rule.State {
	s := state.Data.(dl3057State)

	// Report failures for all bad stages
	for badStage := range s.badStages {
		state = state.AddFailure(rule.CheckFailure{
			Code:     DL3057Meta.Code,
			Severity: DL3057Meta.Severity,
			Message:  DL3057Meta.Message,
			Line:     badStage.line,
		})
	}

	return state
}
