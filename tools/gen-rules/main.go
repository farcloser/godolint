// gen-rules extracts metadata from hadolint Haskell rules and generates Go stubs.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type RuleMetadata struct {
	Code         string
	Severity     string
	Message      string
	SourceFile   string
	Implemented  bool
	HaskellCheck string // Raw Haskell check function for pattern detection
	CanGenerate  bool   // Whether we can auto-generate implementation
}

var (
	// Regex patterns for extracting Haskell rule metadata.
	codePattern     = regexp.MustCompile(`code\s*=\s*"(DL\d+)"`)
	severityPattern = regexp.MustCompile(`severity\s*=\s*(DL\w+)`)
	messagePattern  = regexp.MustCompile(`message\s*=\s*\n?\s*"([^"]*(?:\\\s*\\[^"]*)*)"`)
	ruleTypePattern = regexp.MustCompile(`rule\s*::\s*Rule\s+(.+)`)
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <hadolint-rule-dir>\n", os.Args[0])
		os.Exit(1)
	}

	hadolintRuleDir := os.Args[1]

	// Find all rule files
	pattern := filepath.Join(hadolintRuleDir, "DL*.hs")

	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to glob rules: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No rule files found in %s\n", hadolintRuleDir)
		os.Exit(1)
	}

	// Parse all rules
	var rules []RuleMetadata

	for _, file := range files {
		rule, err := parseRuleFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", file, err)

			continue
		}

		rules = append(rules, rule)
	}

	// Sort by code
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Code < rules[j].Code
	})

	// Check which rules are already implemented (check for _impl.go files)
	for i := range rules {
		implFile := strings.ToLower(rules[i].Code) + "_impl.go"
		if _, err := os.Stat(implFile); err == nil {
			rules[i].Implemented = true
		}
	}

	// Generate output
	fmt.Printf("# Hadolint Rules Status\n\n")
	fmt.Printf("Total rules: %d\n", len(rules))

	implemented := 0

	for _, r := range rules {
		if r.Implemented {
			implemented++
		}
	}

	fmt.Printf("Implemented: %d\n", implemented)
	fmt.Printf("Not implemented: %d\n", len(rules)-implemented)
	fmt.Printf("\n")

	// Always generate metadata files (safe to overwrite)
	metadataGenerated := 0

	for _, rule := range rules {
		if err := generateMetadata(rule); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to generate metadata for %s: %v\n", rule.Code, err)
		} else {
			metadataGenerated++
		}
	}

	fmt.Printf("Generated %d metadata files\n", metadataGenerated)

	// Generate implementations for rules we can auto-generate (never overwrites)
	implGenerated := 0

	for _, rule := range rules {
		if !rule.Implemented && rule.CanGenerate {
			if err := generateImplementation(rule); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to generate implementation for %s: %v\n", rule.Code, err)
			} else {
				implGenerated++
			}
		}
	}

	fmt.Printf("Generated %d working implementations\n", implGenerated)

	// Print summary by status
	fmt.Printf("\n## Rules by Status\n\n")

	fmt.Printf("### Implemented (%d)\n", implemented)

	for _, r := range rules {
		if r.Implemented {
			fmt.Printf("- ✅ %s: %s\n", r.Code, r.Message)
		}
	}

	fmt.Printf("\n### Not Implemented (%d)\n", len(rules)-implemented)

	for _, r := range rules {
		if !r.Implemented {
			fmt.Printf("- ⏳ %s: %s\n", r.Code, r.Message)
		}
	}
}

func parseRuleFile(path string) (RuleMetadata, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return RuleMetadata{}, err
	}

	text := string(content)

	rule := RuleMetadata{
		SourceFile: filepath.Base(path),
	}

	// Extract code
	if match := codePattern.FindStringSubmatch(text); len(match) > 1 {
		rule.Code = match[1]
	} else {
		return rule, errors.New("could not find rule code")
	}

	// Extract severity
	if match := severityPattern.FindStringSubmatch(text); len(match) > 1 {
		rule.Severity = mapSeverity(match[1])
	} else {
		return rule, errors.New("could not find severity")
	}

	// Extract message (handle multiline with \)
	if match := messagePattern.FindStringSubmatch(text); len(match) > 1 {
		msg := match[1]
		// Remove Haskell line continuation
		msg = strings.ReplaceAll(msg, "\\\n", "")
		msg = strings.ReplaceAll(msg, "\\", "")
		msg = strings.TrimSpace(msg)
		rule.Message = msg
	} else {
		return rule, errors.New("could not find message")
	}

	// Extract check function patterns
	// Match patterns like: check (Maintainer _) = False
	//                      check _ = True
	checkPattern := regexp.MustCompile(`(?m)^\s*check\s+(.+?)\s*=`)

	checkMatches := checkPattern.FindAllStringSubmatch(text, -1)
	if len(checkMatches) > 0 {
		// Collect all check patterns
		var checkLines []string

		for _, match := range checkMatches {
			if len(match) > 1 {
				pattern := strings.TrimSpace(match[1])
				// Skip the "where" clause pattern
				if !strings.HasPrefix(pattern, "where") {
					checkLines = append(checkLines, pattern)
				}
			}
		}

		rule.HaskellCheck = strings.Join(checkLines, "\n")
	}

	// Detect if we can auto-generate this rule
	rule.CanGenerate = canGenerateImplementation(rule.HaskellCheck)

	return rule, nil
}

// canGenerateImplementation determines if we can auto-generate Go code from Haskell patterns.
func canGenerateImplementation(haskellCheck string) bool {
	if haskellCheck == "" {
		return false
	}

	lines := strings.Split(haskellCheck, "\n")

	// Only generate for the absolute simplest pattern:
	// Pattern: (InstructionName _) followed by _ (default case)
	// Example: (Maintainer _) and _
	// This means: check specific instruction type, always fail, everything else passes

	if len(lines) != 2 {
		return false
	}

	// First line must be: (InstructionName _) with ONLY underscore, no nested patterns
	simplePattern := regexp.MustCompile(`^\((\w+)\s+_\)$`)
	if !simplePattern.MatchString(strings.TrimSpace(lines[0])) {
		return false
	}

	// Second line must be: _ (default wildcard)
	if strings.TrimSpace(lines[1]) != "_" {
		return false
	}

	return true
}

func mapSeverity(haskellSeverity string) string {
	switch haskellSeverity {
	case "DLErrorC":
		return "rule.Error"
	case "DLWarningC":
		return "rule.Warning"
	case "DLInfoC":
		return "rule.Info"
	case "DLStyleC":
		return "rule.Style"
	default:
		return "rule.Warning" // default
	}
}

// generateMetadata generates metadata-only file (always overwrites).
func generateMetadata(rule RuleMetadata) error {
	tmpl := `// Code generated by go generate; DO NOT EDIT.
// This file is auto-generated from hadolint rules.
package rules

import "github.com/farcloser/godolint/internal/rule"

// {{.Code}}Meta contains metadata for rule {{.Code}}.
// Source: hadolint/src/Hadolint/Rule/{{.SourceFile}}
var {{.Code}}Meta = rule.RuleMeta{
	Code:     "{{.Code}}",
	Severity: {{.Severity}},
	Message:  "{{.Message}}",
}
`

	t, err := template.New("metadata").Parse(tmpl)
	if err != nil {
		return err
	}

	filename := strings.ToLower(rule.Code) + ".go"

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, rule)
}

// generateImplementation generates working implementation from Haskell patterns (only if file doesn't exist).
func generateImplementation(rule RuleMetadata) error {
	filename := strings.ToLower(rule.Code) + "_impl.go"

	// Check if implementation already exists
	if _, err := os.Stat(filename); err == nil {
		return nil // File exists, don't overwrite
	}

	// Generate implementation code based on detected pattern
	implCode := generateCheckFunction(rule)

	tmpl := `package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// {{.Code}} creates a rule from the generated metadata.
// Auto-generated from hadolint Haskell source.
func {{.Code}}() rule.Rule {
	return rule.NewSimpleRule(
		{{.Code}}Meta.Code,
		{{.Code}}Meta.Severity,
		{{.Code}}Meta.Message,
		check{{.Code}},
	)
}

{{.ImplCode}}
`

	t, err := template.New("impl").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	data := struct {
		Code     string
		ImplCode string
	}{
		Code:     rule.Code,
		ImplCode: implCode,
	}

	return t.Execute(f, data)
}

// generateCheckFunction generates Go check function from Haskell patterns.
func generateCheckFunction(rule RuleMetadata) string {
	lines := strings.Split(rule.HaskellCheck, "\n")

	// Pattern 1: Simple instruction type check
	// Example: (Maintainer _) followed by _
	simplePattern := regexp.MustCompile(`\((\w+)\s+[_\w]+\)`)
	if len(lines) == 2 && simplePattern.MatchString(lines[0]) && strings.TrimSpace(lines[1]) == "_" {
		match := simplePattern.FindStringSubmatch(lines[0])
		if len(match) > 1 {
			instructionType := match[1]

			return fmt.Sprintf(`func check%s(instruction syntax.Instruction) bool {
	_, ok := instruction.(*syntax.%s)
	if ok {
		return false // %s instruction found -> fail
	}
	return true
}`, rule.Code, instructionType, instructionType)
		}
	}

	// Fallback - should not reach here if canGenerateImplementation is correct
	return fmt.Sprintf(`func check%s(instruction syntax.Instruction) bool {
	// TODO: Auto-generation failed, implement manually
	return true
}`, rule.Code)
}
