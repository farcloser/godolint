// Package pragma parses hadolint ignore pragmas from Dockerfile comments.
// Ported from Hadolint/Pragma.hs
package pragma

import (
	"regexp"
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// IgnoreDirectives contains parsed ignore pragmas from a Dockerfile.
type IgnoreDirectives struct {
	// LineIgnores maps line numbers to sets of ignored rule codes
	// Key is the line number where the ignore applies (comment line + 1)
	LineIgnores map[int]map[rule.RuleCode]bool
	// GlobalIgnores contains rule codes ignored for the entire file
	GlobalIgnores map[rule.RuleCode]bool
}

// pragmaRegex matches "hadolint ignore=DL3057,DL3018" or "hadolint global ignore=DL3057"
var (
	ignorePragmaRegex       = regexp.MustCompile(`^\s*hadolint\s+ignore\s*=\s*(.+)$`)
	globalIgnorePragmaRegex = regexp.MustCompile(`^\s*hadolint\s+global\s+ignore\s*=\s*(.+)$`)
)

// Parse extracts ignore pragmas from Dockerfile instructions.
// Ported from Hadolint.Pragma module.
func Parse(instructions []syntax.InstructionPos) IgnoreDirectives {
	directives := IgnoreDirectives{
		LineIgnores:   make(map[int]map[rule.RuleCode]bool),
		GlobalIgnores: make(map[rule.RuleCode]bool),
	}

	for _, instr := range instructions {
		comment, ok := instr.Instruction.(*syntax.Comment)
		if !ok {
			continue
		}

		// Check for global ignore pragma
		if codes := parseGlobalIgnorePragma(comment.Text); len(codes) > 0 {
			for _, code := range codes {
				directives.GlobalIgnores[code] = true
			}
			continue
		}

		// Check for line-specific ignore pragma
		if codes := parseIgnorePragma(comment.Text); len(codes) > 0 {
			// Applies to the next line (comment line + 1)
			targetLine := instr.LineNumber + 1
			if directives.LineIgnores[targetLine] == nil {
				directives.LineIgnores[targetLine] = make(map[rule.RuleCode]bool)
			}
			for _, code := range codes {
				directives.LineIgnores[targetLine][code] = true
			}
		}
	}

	return directives
}

// parseIgnorePragma extracts rule codes from "hadolint ignore=DL3057,DL3018" format.
func parseIgnorePragma(text string) []rule.RuleCode {
	matches := ignorePragmaRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil
	}

	return parseRuleList(matches[1])
}

// parseGlobalIgnorePragma extracts rule codes from "hadolint global ignore=DL3057" format.
func parseGlobalIgnorePragma(text string) []rule.RuleCode {
	matches := globalIgnorePragmaRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil
	}

	return parseRuleList(matches[1])
}

// parseRuleList splits comma-separated rule codes and validates format.
// Supports inline comments: "DL3057,DL3018 # some comment"
func parseRuleList(text string) []rule.RuleCode {
	// Strip inline comments (anything after #)
	if idx := strings.Index(text, "#"); idx != -1 {
		text = text[:idx]
	}

	parts := strings.Split(text, ",")
	var codes []rule.RuleCode

	for _, part := range parts {
		code := strings.TrimSpace(part)
		if code == "" {
			continue
		}

		// Validate format: DL followed by 4 digits (DL3057, SC1234, etc)
		if len(code) >= 6 && (code[:2] == "DL" || code[:2] == "SC") {
			codes = append(codes, rule.RuleCode(code))
		}
	}

	return codes
}

// ShouldIgnore returns true if the given failure should be filtered out.
func (d *IgnoreDirectives) ShouldIgnore(failure rule.CheckFailure) bool {
	// Check global ignores
	if d.GlobalIgnores[failure.Code] {
		return true
	}

	// Check line-specific ignores
	if lineIgnores, ok := d.LineIgnores[failure.Line]; ok {
		if lineIgnores[failure.Code] {
			return true
		}
	}

	return false
}