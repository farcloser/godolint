// gen-tests extracts tests from hadolint Haskell test files and generates Go tests.
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

// ErrTestGeneration is the base error for test generation failures.
var ErrTestGeneration = errors.New("test generation error")

// TestCase represents a single test case extracted from hadolint.
type TestCase struct {
	Name       string
	RuleCode   string
	Dockerfile string
	ShouldFail bool // true if ruleCatches, false if ruleCatchesNot
}

// RuleTests contains all test cases for a rule.
type RuleTests struct {
	RuleCode  string
	TestCases []TestCase
	Config    *HadolintConfig
}

// HadolintConfig represents hadolint configuration for test cases.
type HadolintConfig struct {
	LabelSchema  map[string]string // label name -> type (RawText, Email, etc.)
	StrictLabels bool
}

var (
	// Handles escaped quotes within strings: "ADD \"file.zip\" /app/".
	simpleTestPattern = regexp.MustCompile(
		`it\s+"([^"]+)"\s+\$\s+(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+"((?:[^"\\]|\\.)*)"`,
	)

	// Match multi-line end: in ruleCatches "DL3007" $ Text.unlines dockerFile.
	multilineEndPattern = regexp.MustCompile(`in\s+(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+\$\s+Text\.unlines`)
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <hadolint-test-dir>\n", os.Args[0])
		os.Exit(1)
	}

	hadolintTestDir := os.Args[1]

	// Find all test files
	pattern := filepath.Join(hadolintTestDir, "DL*Spec.hs")

	files, err := filepath.Glob(pattern)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to glob tests: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "No test files found in %s\n", hadolintTestDir)
		os.Exit(1)
	}

	// Parse all test files
	allTests := make(map[string][]TestCase)
	allConfigs := make(map[string]*HadolintConfig)
	totalTests := 0

	for _, file := range files {
		tests, configs, err := parseTestFile(file)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", file, err)

			continue
		}

		for ruleCode, cases := range tests {
			allTests[ruleCode] = append(allTests[ruleCode], cases...)
			totalTests += len(cases)
			// Store config if present (last one wins, but should be same per rule)
			if cfg, ok := configs[ruleCode]; ok {
				allConfigs[ruleCode] = cfg
			}
		}
	}

	fmt.Printf("# Hadolint Test Generation\n\n")
	fmt.Printf("Total test files parsed: %d\n", len(files))
	fmt.Printf("Total test cases extracted: %d\n", totalTests)
	fmt.Printf("Rules with tests: %d\n\n", len(allTests))

	// Auto-detect implemented rules by checking for _impl.go files
	implementedRules := []string{}

	for ruleCode := range allTests {
		implFile := strings.ToLower(ruleCode) + "_impl.go"
		if _, err := os.Stat(implFile); err == nil {
			implementedRules = append(implementedRules, ruleCode)
		}
	}

	// Sort for consistent output
	sort.Strings(implementedRules)

	fmt.Printf("Detected %d implemented rules\n\n", len(implementedRules))

	// Generate tests for implemented rules
	generated := 0

	for _, ruleCode := range implementedRules {
		cases, ok := allTests[ruleCode]
		if !ok {
			fmt.Printf("Warning: No tests found for implemented rule %s\n", ruleCode)

			continue
		}

		config := allConfigs[ruleCode] // May be nil if no config found

		if err := generateTestFile(ruleCode, cases, config); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error generating tests for %s: %v\n", ruleCode, err)
		} else {
			fmt.Printf("Generated tests for %s (%d cases)\n", ruleCode, len(cases))

			generated++
		}
	}

	fmt.Printf("\nGenerated test files: %d\n", generated)
}

func parseTestFile(path string) (map[string][]TestCase, map[string]*HadolintConfig, error) {
	//nolint:gosec
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrTestGeneration, err)
	}

	text := string(content)
	tests := make(map[string][]TestCase)
	configs := make(map[string]*HadolintConfig)

	// Parse config for this file
	config := parseHadolintConfig(text)

	// Parse simple single-line tests
	matches := simpleTestPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) < 5 {
			continue
		}

		testCase := TestCase{
			Name:       match[1],
			ShouldFail: match[2] == "ruleCatches",
			RuleCode:   match[3],
			Dockerfile: unescapeHaskellString(match[4]),
		}

		tests[testCase.RuleCode] = append(tests[testCase.RuleCode], testCase)

		if config != nil {
			configs[testCase.RuleCode] = config
		}
	}

	// Parse do-block tests
	doBlockTests := parseDoBlockTests(text)
	for _, testCase := range doBlockTests {
		tests[testCase.RuleCode] = append(tests[testCase.RuleCode], testCase)

		if config != nil {
			configs[testCase.RuleCode] = config
		}
	}

	// Parse multi-line tests
	multilineTests := parseMultilineTests(text)
	for _, testCase := range multilineTests {
		tests[testCase.RuleCode] = append(tests[testCase.RuleCode], testCase)

		if config != nil {
			configs[testCase.RuleCode] = config
		}
	}

	return tests, configs, nil
}

// unescapeHaskellString converts Haskell escape sequences to actual characters.
// Implements complete Haskell single-character escape sequences per Haskell 98 spec.
func unescapeHaskellString(s string) string {
	// Process escape sequences using regex to handle them correctly
	escapePattern := regexp.MustCompile(`\\(.)`)
	result := escapePattern.ReplaceAllStringFunc(s, func(match string) string {
		// match is like "\n" or "\t" or "\\" or "\""
		if len(match) != 2 {
			return match
		}

		switch match[1] {
		case '0':
			return "\x00" // null
		case 'a':
			return "\x07" // alert (bell)
		case 'b':
			return "\x08" // backspace
		case 'f':
			return "\x0C" // form feed
		case 'n':
			return "\n" // newline
		case 'r':
			return "\r" // carriage return
		case 't':
			return "\t" // horizontal tab
		case 'v':
			return "\x0B" // vertical tab
		case '"':
			return `"` // double quote
		case '\'':
			return "'" // single quote
		case '\\':
			return `\` // backslash
		case '&':
			return "" // empty string (Haskell string gap)
		default:
			// Unknown escape sequence, keep as-is
			return match
		}
	})

	return result
}

// Example: let ?config = def { labelSchema = Map.fromList [("foo", Rule.RawText)], strictLabels = True }.
func parseHadolintConfig(text string) *HadolintConfig {
	config := &HadolintConfig{
		LabelSchema: make(map[string]string),
	}

	// Match: let ?config = def { ... }
	configPattern := regexp.MustCompile(`let\s+\?config\s*=\s*def\s*\{([^}]+)}`)

	configMatch := configPattern.FindStringSubmatch(text)
	if configMatch == nil {
		return nil // No config found
	}

	configBody := configMatch[1]

	// Extract labelSchema: Map.fromList [("label1", Rule.Type1), ...]
	labelSchemaPattern := regexp.MustCompile(`labelSchema\s*=\s*Map\.fromList\s*\[(.*?)]`)

	labelMatch := labelSchemaPattern.FindStringSubmatch(configBody)
	if labelMatch != nil {
		// Parse label entries: ("label", Rule.Type) or ("label", Type)
		entryPattern := regexp.MustCompile(`\("([^"]+)",\s*(?:Rule\.)?(\w+)\)`)

		entries := entryPattern.FindAllStringSubmatch(labelMatch[1], -1)
		for _, entry := range entries {
			if len(entry) >= 3 {
				labelName := entry[1]
				labelType := entry[2]
				config.LabelSchema[labelName] = labelType
			}
		}
	}

	// Extract strictLabels
	strictPattern := regexp.MustCompile(`strictLabels\s*=\s*(True|False)`)

	strictMatch := strictPattern.FindStringSubmatch(configBody)
	if strictMatch != nil {
		config.StrictLabels = strictMatch[1] == "True"
	}

	// Return nil if no config was actually found
	if len(config.LabelSchema) == 0 && !config.StrictLabels {
		return nil
	}

	return config
}

func parseDoBlockTests(text string) []TestCase {
	var tests []TestCase

	lines := strings.Split(text, "\n")

	// Pattern to match: it "test name" $ do
	doPattern := regexp.MustCompile(`it\s+"([^"]+)"\s+\$\s+do\s*$`)
	// Pattern to match assertions: ruleCatches "DL3048" "LABEL ..."
	// or: onBuildRuleCatches "DL3048" "LABEL ..."
	assertionPattern := regexp.MustCompile(
		`\s*(onBuild)?(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+"((?:[^"\\]|\\.)*)"`,
	)

	for idx := range lines {
		line := lines[idx]

		// Look for "it 'name' $ do" line
		doMatch := doPattern.FindStringSubmatch(line)
		if doMatch == nil {
			continue
		}

		testName := doMatch[1]
		baseIndent := len(line) - len(strings.TrimLeft(line, " "))

		// Collect all assertions in this do block
		assertionCount := 0

		for j := idx + 1; j < len(lines); j++ {
			currentLine := lines[j]

			// Empty lines are ok
			if strings.TrimSpace(currentLine) == "" {
				continue
			}

			// Check if we've exited the do block (decreased indentation to base level or less)
			currentIndent := len(currentLine) - len(strings.TrimLeft(currentLine, " "))
			if currentIndent <= baseIndent {
				break
			}

			// Look for assertion pattern
			assertMatch := assertionPattern.FindStringSubmatch(currentLine)
			if assertMatch != nil {
				assertionCount++
				shouldFail := assertMatch[2] == "ruleCatches"
				ruleCode := assertMatch[3]
				dockerfile := unescapeHaskellString(assertMatch[4])

				// Create unique test name if multiple assertions
				name := testName
				if assertionCount > 1 {
					name = fmt.Sprintf("%s (%d)", testName, assertionCount)
				}

				tests = append(tests, TestCase{
					Name:       name,
					RuleCode:   ruleCode,
					Dockerfile: dockerfile,
					ShouldFail: shouldFail,
				})
			}
		}
	}

	return tests
}

//revive:disable:max-control-nesting // yolo!
func parseMultilineTests(text string) []TestCase {
	var tests []TestCase

	// This parser handles multiple let-binding formats:
	// Format 1 (dockerFile list):
	//   it "name" $
	//     let dockerFile = [...]
	//      in do
	//           ruleCatches "DL3009" $ Text.unlines dockerFile
	//
	// Format 2 (line string):
	//   it "name" $
	//     let line = "RUN ..."
	//     in do
	//         ruleCatches "DL3019" line
	//
	// Format 3 (in do on same line as assertion):
	//   it "name" $
	//     let dockerFile = [...]
	//      in do ruleCatches "DL3047" $ Text.unlines dockerFile

	lines := strings.Split(text, "\n")
	idx := 0

	// Pattern to match: it "test name" $
	itPattern := regexp.MustCompile(`it\s+"([^"]+)"\s+\$\s*$`)
	// Pattern to match: let dockerFile = or let line =
	letDockerFilePattern := regexp.MustCompile(`^\s*let\s+dockerFile\s*=`)
	letLinePattern := regexp.MustCompile(`^\s*let\s+line\s*=\s*"((?:[^"\\]|\\.)*)"`)
	// Pattern to match: in do (standalone)
	inDoPattern := regexp.MustCompile(`^\s*in\s+do\s*$`)
	// Pattern to match: in do ruleCatches... (inline)
	inDoInlinePattern := regexp.MustCompile(`^\s*in\s+do\s+(onBuild)?(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+`)
	// Pattern to match assertions: ruleCatches "DL3009" $ Text.unlines dockerFile
	// or: ruleCatches "DL3019" line
	assertionPattern := regexp.MustCompile(
		`\s*(onBuild)?(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+(\$\s+Text\.unlines\s+dockerFile|line)`,
	)

	for idx < len(lines) {
		line := lines[idx]

		// Look for "it 'name' $" line
		itMatch := itPattern.FindStringSubmatch(line)
		if itMatch == nil {
			idx++

			continue
		}

		testName := itMatch[1]

		// Check if next line has a let binding
		if idx+1 >= len(lines) {
			idx++

			continue
		}

		nextLine := lines[idx+1]

		// Check for "let line = ..." format
		lineMatch := letLinePattern.FindStringSubmatch(nextLine)
		if lineMatch != nil {
			// Single line dockerfile
			dockerfile := unescapeHaskellString(lineMatch[1])

			// Find "in do" and parse assertions
			foundInDo := false
			inDoLineIndex := -1

			for j := idx + 2; j < len(lines) && j < idx+20; j++ {
				currentLine := lines[j]

				// Look for "in do" pattern (standalone or inline)
				if inDoPattern.MatchString(currentLine) {
					foundInDo = true
					inDoLineIndex = j

					break
				}

				// Check for inline "in do ruleCatches..."
				inlineMatch := inDoInlinePattern.FindStringSubmatch(currentLine)
				if inlineMatch != nil {
					foundInDo = true
					inDoLineIndex = j

					// Parse the inline assertion immediately
					shouldFail := inlineMatch[2] == "ruleCatches"
					ruleCode := inlineMatch[3]

					tests = append(tests, TestCase{
						Name:       testName,
						RuleCode:   ruleCode,
						Dockerfile: dockerfile,
						ShouldFail: shouldFail,
					})

					// Continue parsing subsequent lines for more assertions
					break
				}
			}

			if foundInDo && inDoLineIndex >= 0 {
				baseIndent := len(lines[inDoLineIndex]) - len(strings.TrimLeft(lines[inDoLineIndex], " "))
				assertionCount := 0

				for j := inDoLineIndex + 1; j < len(lines) && j < inDoLineIndex+20; j++ {
					currentLine := lines[j]

					if strings.TrimSpace(currentLine) == "" {
						continue
					}

					currentIndent := len(currentLine) - len(strings.TrimLeft(currentLine, " "))
					if currentIndent <= baseIndent {
						idx = j - 1

						break
					}

					assertMatch := assertionPattern.FindStringSubmatch(currentLine)
					if assertMatch != nil {
						assertionCount++
						shouldFail := assertMatch[2] == "ruleCatches"
						ruleCode := assertMatch[3]

						name := testName
						if assertionCount > 1 {
							name = fmt.Sprintf("%s (%d)", testName, assertionCount)
						}

						tests = append(tests, TestCase{
							Name:       name,
							RuleCode:   ruleCode,
							Dockerfile: dockerfile,
							ShouldFail: shouldFail,
						})
					}
				}
			}

			idx++

			continue
		}

		// Check for "let dockerFile = [...]" format
		if !letDockerFilePattern.MatchString(nextLine) {
			idx++

			continue
		}

		// Collect dockerfile lines from the let binding
		var dockerfileLines []string

		foundInDo := false
		inDoLineIndex := -1

		for j := idx + 1; j < len(lines) && j < idx+50; j++ {
			currentLine := lines[j]

			// Check for terminating patterns FIRST (before extracting quotes)
			// to avoid including rule codes in Dockerfile content

			// Look for "in do" pattern (standalone)
			if inDoPattern.MatchString(currentLine) {
				foundInDo = true
				inDoLineIndex = j

				break
			}

			// Check for inline "in do ruleCatches..."
			inlineMatch := inDoInlinePattern.FindStringSubmatch(currentLine)
			if inlineMatch != nil {
				foundInDo = true
				inDoLineIndex = j

				break
			}

			// Also support the old format: in ruleCatches (without do-block)
			endMatch := multilineEndPattern.FindStringSubmatch(currentLine)
			if endMatch != nil {
				shouldFail := endMatch[1] == "ruleCatches"
				ruleCode := endMatch[2]

				if len(dockerfileLines) > 0 {
					tests = append(tests, TestCase{
						Name:       testName,
						RuleCode:   ruleCode,
						Dockerfile: strings.Join(dockerfileLines, "\n"),
						ShouldFail: shouldFail,
					})
				}

				idx = j

				break
			}

			// Extract dockerfile lines from quoted strings (only if not a terminating pattern)
			if strings.Contains(currentLine, "\"") {
				// Extract quoted strings, handling escaped quotes
				quotedStrings := regexp.MustCompile(`"((?:[^"\\]|\\.)*)"`).FindAllStringSubmatch(currentLine, -1)
				for _, match := range quotedStrings {
					if len(match) > 1 && match[1] != "" {
						// Unescape Haskell string escape sequences
						unescaped := unescapeHaskellString(match[1])
						dockerfileLines = append(dockerfileLines, unescaped)
					}
				}
			}
		}

		// If we found "in do", parse the assertions in the do-block
		if foundInDo && len(dockerfileLines) > 0 {
			dockerfile := strings.Join(dockerfileLines, "\n")

			// Check if it's inline format: "in do ruleCatches..."
			inlineMatch := inDoInlinePattern.FindStringSubmatch(lines[inDoLineIndex])
			if inlineMatch != nil {
				// Parse the inline assertion
				shouldFail := inlineMatch[2] == "ruleCatches"
				ruleCode := inlineMatch[3]

				tests = append(tests, TestCase{
					Name:       testName,
					RuleCode:   ruleCode,
					Dockerfile: dockerfile,
					ShouldFail: shouldFail,
				})

				// Continue parsing subsequent lines for more assertions
				baseIndent := len(lines[inDoLineIndex]) - len(strings.TrimLeft(lines[inDoLineIndex], " "))
				assertionCount := 1

				for j := inDoLineIndex + 1; j < len(lines) && j < inDoLineIndex+20; j++ {
					currentLine := lines[j]

					if strings.TrimSpace(currentLine) == "" {
						continue
					}

					currentIndent := len(currentLine) - len(strings.TrimLeft(currentLine, " "))
					if currentIndent <= baseIndent {
						idx = j - 1

						break
					}

					assertMatch := assertionPattern.FindStringSubmatch(currentLine)
					if assertMatch != nil {
						assertionCount++
						shouldFail := assertMatch[2] == "ruleCatches"
						ruleCode := assertMatch[3]

						name := fmt.Sprintf("%s (%d)", testName, assertionCount)

						tests = append(tests, TestCase{
							Name:       name,
							RuleCode:   ruleCode,
							Dockerfile: dockerfile,
							ShouldFail: shouldFail,
						})
					}
				}
			} else {
				// Standalone "in do" - parse assertions from next line
				baseIndent := len(lines[inDoLineIndex]) - len(strings.TrimLeft(lines[inDoLineIndex], " "))
				assertionCount := 0

				for j := inDoLineIndex + 1; j < len(lines) && j < inDoLineIndex+20; j++ {
					currentLine := lines[j]

					// Empty lines are ok
					if strings.TrimSpace(currentLine) == "" {
						continue
					}

					// Check if we've exited the do block (decreased indentation to base level or less)
					currentIndent := len(currentLine) - len(strings.TrimLeft(currentLine, " "))
					if currentIndent <= baseIndent {
						idx = j - 1

						break
					}

					// Look for assertion pattern
					assertMatch := assertionPattern.FindStringSubmatch(currentLine)
					if assertMatch != nil {
						assertionCount++
						shouldFail := assertMatch[2] == "ruleCatches"
						ruleCode := assertMatch[3]

						// Create unique test name if multiple assertions
						name := testName
						if assertionCount > 1 {
							name = fmt.Sprintf("%s (%d)", testName, assertionCount)
						}

						tests = append(tests, TestCase{
							Name:       name,
							RuleCode:   ruleCode,
							Dockerfile: dockerfile,
							ShouldFail: shouldFail,
						})
					}
				}
			}
		}

		idx++
	}

	return tests
}

func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "` + \"`\" + `")
}

// goLabelType maps Haskell LabelType names to Go LabelType names.
func goLabelType(haskellType string) string {
	// Map Haskell type names to Go type names (handle acronym casing differences)
	typeMap := map[string]string{
		"RawText": "RawText",
		"Email":   "Email",
		"Url":     "URL",     // Haskell: Url, Go: URL
		"Rfc3339": "RFC3339", // Haskell: Rfc3339, Go: RFC3339
		"Spdx":    "SPDX",    // Haskell: Spdx, Go: SPDX
		"GitHash": "GitHash",
		"SemVer":  "SemVer",
	}
	if goType, ok := typeMap[haskellType]; ok {
		return goType
	}

	return haskellType // fallback
}

func generateTestFile(ruleCode string, cases []TestCase, config *HadolintConfig) error {
	// Sort test cases by name for consistent output
	sort.Slice(cases, func(i, j int) bool {
		return cases[i].Name < cases[j].Name
	})

	tmpl := `package rules_test

import (
	"testing"

{{if .Config}}	"github.com/farcloser/godolint/internal/config"
{{end}}	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for {{.RuleCode}} ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/{{.RuleCode}}Spec.hs
//
// To regenerate: go generate ./internal/rules

func Test{{.RuleCode}}(t *testing.T) {
	t.Parallel()

{{if .Config}}	cfg := &config.Config{
{{if .Config.LabelSchema}}		LabelSchema: map[string]config.LabelType{
{{range $key, $val := .Config.LabelSchema}}			"{{$key}}": config.LabelType{{goLabelType $val}},
{{end}}		},
{{end}}{{if .Config.StrictLabels}}		StrictLabels: true,
{{end}}	}
	allRules := []rule.Rule{
		rules.{{.RuleCode}}WithConfig(cfg),
	}
{{else}}	allRules := []rule.Rule{
		rules.{{.RuleCode}}(),
	}
{{end}}{{range .TestCases}}
	t.Run(
		{{.Name | printf "%q"}},
		func(t *testing.T) {
			t.Parallel()

			dockerfile := ` + "`" + `{{.Dockerfile | escapeBackticks}}` + "`" + `
			violations := testutils.LintDockerfile(dockerfile, allRules)
{{if .ShouldFail}}
			testutils.AssertContainsViolation(t, violations, "{{.RuleCode}}"){{else}}
			testutils.AssertNoViolation(t, violations, "{{.RuleCode}}"){{end}}
		},
	)
{{end}}}
`

	funcMap := template.FuncMap{
		"escapeBackticks": escapeBackticks,
		"goLabelType":     goLabelType,
	}

	tpl, err := template.New("test").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTestGeneration, err)
	}

	filename := strings.ToLower(ruleCode) + "_test.go"

	//nolint:gosec
	outputFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTestGeneration, err)
	}

	defer func() {
		_ = outputFile.Close()
	}()

	data := RuleTests{
		RuleCode:  ruleCode,
		TestCases: cases,
		Config:    config,
	}

	err = tpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTestGeneration, err)
	}

	return nil
}
