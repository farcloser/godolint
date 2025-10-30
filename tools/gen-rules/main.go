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
	Code          string
	Severity      string
	Message       string
	RequiresShell bool
	SourceFile    string
	Implemented   bool
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

	// Check which rules are already implemented
	for i := range rules {
		goFile := fmt.Sprintf("%s.go", strings.ToLower(rules[i].Code))
		if _, err := os.Stat(goFile); err == nil {
			rules[i].Implemented = true
		}
	}

	// Generate output
	fmt.Printf("# Hadolint Rules Status\n\n")
	fmt.Printf("Total rules: %d\n", len(rules))

	implemented := 0
	requiresShell := 0
	canImplement := 0

	for _, r := range rules {
		if r.Implemented {
			implemented++
		}

		if r.RequiresShell {
			requiresShell++
		} else {
			canImplement++
		}
	}

	fmt.Printf("Implemented: %d\n", implemented)
	fmt.Printf("Requires shell parsing: %d\n", requiresShell)
	fmt.Printf("Can implement without shell: %d\n", canImplement)
	fmt.Printf("\n")

	// Generate stubs for unimplemented rules
	stubsGenerated := 0

	for _, rule := range rules {
		if !rule.Implemented && !rule.RequiresShell {
			if err := generateStub(rule); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to generate stub for %s: %v\n", rule.Code, err)
			} else {
				stubsGenerated++
			}
		}
	}

	fmt.Printf("Generated %d new rule stubs\n", stubsGenerated)

	// Print summary by status
	fmt.Printf("\n## Rules by Status\n\n")

	fmt.Printf("### Implemented (%d)\n", implemented)

	for _, r := range rules {
		if r.Implemented {
			fmt.Printf("- ✅ %s: %s\n", r.Code, r.Message)
		}
	}

	fmt.Printf("\n### Ready to Implement (%d - no shell parsing required)\n", canImplement-implemented)

	for _, r := range rules {
		if !r.Implemented && !r.RequiresShell {
			fmt.Printf("- ⏳ %s: %s\n", r.Code, r.Message)
		}
	}

	fmt.Printf("\n### Requires Shell Parsing (%d)\n", requiresShell)

	for _, r := range rules {
		if r.RequiresShell {
			impl := "❌"
			if r.Implemented {
				impl = "✅"
			}

			fmt.Printf("- %s %s: %s\n", impl, r.Code, r.Message)
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

	// Check if requires shell parsing
	if match := ruleTypePattern.FindStringSubmatch(text); len(match) > 1 {
		ruleType := match[1]
		rule.RequiresShell = strings.Contains(ruleType, "ParsedShell")
	}

	return rule, nil
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

func generateStub(rule RuleMetadata) error {
	tmpl := `package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// {{.Code}} - {{.Message}}
//
// Ported from Hadolint.Rule.{{.Code}}
// Source: {{.SourceFile}}
func {{.Code}}() rule.Rule {
	return rule.NewSimpleRule(
		"{{.Code}}",
		{{.Severity}},
		"{{.Message}}",
		check{{.Code}},
	)
}

func check{{.Code}}(instruction syntax.Instruction) bool {
	// TODO: Port check logic from hadolint/src/Hadolint/Rule/{{.SourceFile}}
	// See: hadolint/src/Hadolint/Rule/{{.SourceFile}} for implementation

	// Placeholder: allow all instructions for now
	return true
}
`

	t, err := template.New("rule").Parse(tmpl)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.go", strings.ToLower(rule.Code))

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, rule)
}
