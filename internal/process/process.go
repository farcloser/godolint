// Package process orchestrates running hadolint rules over a parsed Dockerfile AST.
// Ported from Hadolint/Process.hs
package process

import (
	"github.com/farcloser/godolint/internal/pragma"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Processor runs rules against a Dockerfile AST and collects violations.
type Processor struct {
	rules                 []rule.Rule
	disableIgnorePragmas bool
}

// NewProcessor creates a new processor with the given rules.
func NewProcessor(rules []rule.Rule) *Processor {
	return &Processor{
		rules:                 rules,
		disableIgnorePragmas: false,
	}
}

// WithDisableIgnorePragmas configures whether to disable inline ignore pragma processing.
func (p *Processor) WithDisableIgnorePragmas(disable bool) *Processor {
	p.disableIgnorePragmas = disable
	return p
}

// Run processes a Dockerfile AST and returns all rule violations found.
// Uses fold-style accumulation with state for each rule.
// Ported from Hadolint's Rule fold pattern.
func (p *Processor) Run(instructions []syntax.InstructionPos) []rule.CheckFailure {
	allFailures := []rule.CheckFailure{}

	// For each rule, fold over all instructions with state
	for _, currentRule := range p.rules {
		state := currentRule.InitialState()

		// Thread state through each instruction check
		for _, instrPos := range instructions {
			state = currentRule.Check(instrPos.LineNumber, state, instrPos.Instruction)
		}

		// Finalize the state (some rules add failures only at the end)
		state = currentRule.Finalize(state)

		// Collect failures from final state
		allFailures = append(allFailures, state.Failures...)
	}

	// Filter out failures with Ignore severity (like hadolint's DLIgnoreC filter)
	// Ported from Hadolint/Lint.hs:88 - severity /= DLIgnoreC
	allFailures = filterIgnoreSeverity(allFailures)

	// Filter out ignored failures based on inline pragmas
	if !p.disableIgnorePragmas {
		directives := pragma.Parse(instructions)
		allFailures = filterIgnored(allFailures, directives)
	}

	return allFailures
}

// filterIgnoreSeverity removes failures with Ignore severity.
// Matches hadolint's behavior where DLIgnoreC severity rules are filtered out.
func filterIgnoreSeverity(failures []rule.CheckFailure) []rule.CheckFailure {
	filtered := []rule.CheckFailure{}
	for _, failure := range failures {
		if failure.Severity != rule.Ignore {
			filtered = append(filtered, failure)
		}
	}
	return filtered
}

// filterIgnored removes failures that are suppressed by ignore pragmas.
func filterIgnored(failures []rule.CheckFailure, directives pragma.IgnoreDirectives) []rule.CheckFailure {
	filtered := []rule.CheckFailure{}
	for _, failure := range failures {
		if !directives.ShouldIgnore(failure) {
			filtered = append(filtered, failure)
		}
	}
	return filtered
}
