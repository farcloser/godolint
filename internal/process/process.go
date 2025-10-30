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
// Each instruction is checked against all registered rules.
func (p *Processor) Run(instructions []syntax.InstructionPos) []rule.CheckFailure {
	var failures []rule.CheckFailure

	for _, instrPos := range instructions {
		for _, r := range p.rules {
			if !r.Check(instrPos.Instruction) {
				failures = append(failures, rule.CheckFailure{
					Code:     r.Code(),
					Severity: r.Severity(),
					Message:  r.Message(),
					Line:     instrPos.LineNumber,
				})
			}
		}
	}

	return failures
}
