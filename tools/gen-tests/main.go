// gen-tests extracts tests from hadolint Haskell test files and generates Go tests.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type TestCase struct {
	Name       string
	RuleCode   string
	Dockerfile string
	ShouldFail bool // true if ruleCatches, false if ruleCatchesNot
}

type RuleTests struct {
	RuleCode  string
	TestCases []TestCase
}

var (
	// Match: it "test name" $ ruleCatches "DL3007" "FROM debian:latest"
	// Handles escaped quotes within strings: "ADD \"file.zip\" /app/"
	simpleTestPattern = regexp.MustCompile(`it\s+"([^"]+)"\s+\$\s+(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+"((?:[^"\\]|\\.)*)"`)

	// Match multi-line end: in ruleCatches "DL3007" $ Text.unlines dockerFile
	multilineEndPattern = regexp.MustCompile(`in\s+(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+\$\s+Text\.unlines`)
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <hadolint-test-dir>\n", os.Args[0])
		os.Exit(1)
	}

	hadolintTestDir := os.Args[1]

	// Find all test files
	pattern := filepath.Join(hadolintTestDir, "DL*Spec.hs")
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to glob tests: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No test files found in %s\n", hadolintTestDir)
		os.Exit(1)
	}

	// Parse all test files
	allTests := make(map[string][]TestCase)
	totalTests := 0

	for _, file := range files {
		tests, err := parseTestFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		for ruleCode, cases := range tests {
			allTests[ruleCode] = append(allTests[ruleCode], cases...)
			totalTests += len(cases)
		}
	}

	fmt.Printf("# Hadolint Test Generation\n\n")
	fmt.Printf("Total test files parsed: %d\n", len(files))
	fmt.Printf("Total test cases extracted: %d\n", totalTests)
	fmt.Printf("Rules with tests: %d\n\n", len(allTests))

	// Auto-detect implemented rules by checking for _impl.go files
	implementedRules := []string{}
	for ruleCode := range allTests {
		implFile := fmt.Sprintf("%s_impl.go", strings.ToLower(ruleCode))
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

		if err := generateTestFile(ruleCode, cases); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating tests for %s: %v\n", ruleCode, err)
		} else {
			fmt.Printf("Generated tests for %s (%d cases)\n", ruleCode, len(cases))
			generated++
		}
	}

	fmt.Printf("\nGenerated test files: %d\n", generated)
}

func parseTestFile(path string) (map[string][]TestCase, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := string(content)
	tests := make(map[string][]TestCase)

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
	}

	// Parse do-block tests
	doBlockTests := parseDoBlockTests(text)
	for _, testCase := range doBlockTests {
		tests[testCase.RuleCode] = append(tests[testCase.RuleCode], testCase)
	}

	// Parse multi-line tests
	multilineTests := parseMultilineTests(text)
	for _, testCase := range multilineTests {
		tests[testCase.RuleCode] = append(tests[testCase.RuleCode], testCase)
	}

	return tests, nil
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

func parseDoBlockTests(text string) []TestCase {
	var tests []TestCase

	lines := strings.Split(text, "\n")

	// Pattern to match: it "test name" $ do
	doPattern := regexp.MustCompile(`it\s+"([^"]+)"\s+\$\s+do\s*$`)
	// Pattern to match assertions: ruleCatches "DL3048" "LABEL ..."
	// or: onBuildRuleCatches "DL3048" "LABEL ..."
	assertionPattern := regexp.MustCompile(`\s*(onBuild)?(ruleCatches|ruleCatchesNot)\s+"(DL\d+)"\s+"((?:[^"\\]|\\.)*)"`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Look for "it 'name' $ do" line
		doMatch := doPattern.FindStringSubmatch(line)
		if doMatch == nil {
			continue
		}

		testName := doMatch[1]
		baseIndent := len(line) - len(strings.TrimLeft(line, " "))

		// Collect all assertions in this do block
		assertionCount := 0
		for j := i + 1; j < len(lines); j++ {
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

func parseMultilineTests(text string) []TestCase {
	var tests []TestCase

	// This is a simplified parser for multi-line tests
	// Find blocks that match: it "name" $ (newline) let dockerFile = [...] in ruleCatches ...

	lines := strings.Split(text, "\n")
	i := 0

	// Pattern to match: it "test name" $
	itPattern := regexp.MustCompile(`it\s+"([^"]+)"\s+\$\s*$`)
	// Pattern to match: let dockerFile =
	letPattern := regexp.MustCompile(`^\s*let\s+dockerFile\s*=`)

	for i < len(lines) {
		line := lines[i]

		// Look for "it 'name' $" line
		itMatch := itPattern.FindStringSubmatch(line)
		if itMatch == nil {
			i++
			continue
		}

		testName := itMatch[1]

		// Check if next line starts with "let dockerFile ="
		if i+1 < len(lines) && !letPattern.MatchString(lines[i+1]) {
			i++
			continue
		}

		// Collect lines until we find the closing "in ruleCatches/ruleCatchesNot"
		var dockerfileLines []string
		foundEnd := false
		ruleCode := ""
		shouldFail := false

		for j := i; j < len(lines) && j < i+50; j++ {
			currentLine := lines[j]

			// Look for dockerfile lines in the list
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

			// Look for end pattern
			endMatch := multilineEndPattern.FindStringSubmatch(currentLine)
			if endMatch != nil {
				shouldFail = endMatch[1] == "ruleCatches"
				ruleCode = endMatch[2]
				foundEnd = true
				i = j
				break
			}
		}

		if foundEnd && len(dockerfileLines) > 0 {
			tests = append(tests, TestCase{
				Name:       testName,
				RuleCode:   ruleCode,
				Dockerfile: strings.Join(dockerfileLines, "\n"),
				ShouldFail: shouldFail,
			})
		}

		i++
	}

	return tests
}

func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "` + \"`\" + `")
}

func generateTestFile(ruleCode string, cases []TestCase) error {
	// Sort test cases by name for consistent output
	sort.Slice(cases, func(i, j int) bool {
		return cases[i].Name < cases[j].Name
	})

	tmpl := `package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for {{.RuleCode}} ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/{{.RuleCode}}Spec.hs
//
// To regenerate: go generate ./internal/rules

func Test{{.RuleCode}}(t *testing.T) {
	allRules := []rule.Rule{ {{.RuleCode}}() }

{{range .TestCases}}
	t.Run({{.Name | printf "%q"}}, func(t *testing.T) {
		dockerfile := ` + "`" + `{{.Dockerfile | escapeBackticks}}` + "`" + `
		violations := LintDockerfile(dockerfile, allRules)
{{if .ShouldFail}}
		AssertContainsViolation(t, violations, "{{.RuleCode}}")
{{else}}
		AssertNoViolation(t, violations, "{{.RuleCode}}")
{{end}}
	})
{{end}}
}
`

	funcMap := template.FuncMap{
		"escapeBackticks": escapeBackticks,
	}

	t, err := template.New("test").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_test.go", strings.ToLower(ruleCode))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		RuleCode  string
		TestCases []TestCase
	}{
		RuleCode:  ruleCode,
		TestCases: cases,
	}

	return t.Execute(f, data)
}
