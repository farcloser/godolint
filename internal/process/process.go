// Package process orchestrates running hadolint rules over a parsed Dockerfile AST.
// Ported from Hadolint/Process.hs
package process

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Processor runs rules against a Dockerfile AST and collects violations.
type Processor struct {
	rules []rule.Rule
}

// NewProcessor creates a new processor with the given rules.
func NewProcessor(rules []rule.Rule) *Processor {
	return &Processor{
		rules: rules,
	}
}

// Run processes a Dockerfile AST and returns all rule violations found.
// Uses fold-style accumulation with state for each rule.
// Ported from Hadolint's Rule fold pattern.
func (p *Processor) Run(instructions []syntax.InstructionPos) []rule.CheckFailure {
	var allFailures []rule.CheckFailure

	// For each rule, fold over all instructions with state
	for _, r := range p.rules {
		state := r.InitialState()

		// Thread state through each instruction check
		for _, instrPos := range instructions {
			state = r.Check(instrPos.LineNumber, state, instrPos.Instruction)
		}

		// Finalize the state (some rules add failures only at the end)
		state = r.Finalize(state)

		// Collect failures from final state
		allFailures = append(allFailures, state.Failures...)
	}

	return allFailures
}
