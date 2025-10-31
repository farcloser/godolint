package rules

import (
	"strings"
	"unicode"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3033 checks for yum install without version pinning.
func DL3033() rule.Rule {
	return rule.NewSimpleRule(
		DL3033Meta.Code,
		DL3033Meta.Severity,
		DL3033Meta.Message,
		checkDL3033,
	)
}

func checkDL3033(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all yum packages
	packages := getYumPackages(parsed)
	for _, pkg := range packages {
		if !isYumPackageVersionFixed(pkg) {
			return false
		}
	}

	// Check all yum modules
	modules := getYumModules(parsed)
	for _, mod := range modules {
		if !isYumModuleVersionFixed(mod) {
			return false
		}
	}

	return true
}

func getYumPackages(parsed *shell.ParsedShell) []string {
	var packages []string

	for _, cmd := range parsed.PresentCommands {
		// Skip yum module commands
		if shell.CmdHasArgs("yum", []string{"module"}, cmd) {
			continue
		}

		// Get packages from install filter
		packages = append(packages, yumInstallFilter(cmd)...)
	}

	return packages
}

func getYumModules(parsed *shell.ParsedShell) []string {
	var modules []string

	for _, cmd := range parsed.PresentCommands {
		// Only yum module commands
		if !shell.CmdHasArgs("yum", []string{"module"}, cmd) {
			continue
		}

		// Get modules from install filter
		modules = append(modules, yumInstallFilter(cmd)...)
	}

	return modules
}

func yumInstallFilter(cmd shell.Command) []string {
	// Must be yum install
	if !shell.CmdHasArgs("yum", []string{"install"}, cmd) {
		return nil
	}

	args := shell.GetArgsNoFlags(cmd)

	var packages []string

	for _, arg := range args {
		if arg != "install" && arg != "module" {
			packages = append(packages, arg)
		}
	}

	return packages
}

func isYumPackageVersionFixed(pkg string) bool {
	// .rpm files always have version
	if strings.HasSuffix(pkg, ".rpm") {
		return true
	}

	// Must have at least one dash
	parts := strings.Split(pkg, "-")
	if len(parts) <= 1 {
		return false
	}

	// Check if parts after first dash look like version
	versionParts := parts[1:]

	return isVersionLike(versionParts)
}

func isVersionLike(parts []string) bool {
	if len(parts) == 0 {
		return false
	}

	allValid := true
	hasDigitStart := false

	for _, part := range parts {
		if !isValidVersionPart(part) {
			allValid = false

			break
		}

		if len(part) > 0 && unicode.IsDigit(rune(part[0])) {
			hasDigitStart = true
		}
	}

	return allValid && hasDigitStart
}

func isValidVersionPart(part string) bool {
	if len(part) == 0 {
		return false
	}

	for _, ch := range part {
		if !isVersionChar(ch) {
			return false
		}
	}

	return true
}

func isVersionChar(ch rune) bool {
	return unicode.IsDigit(ch) ||
		unicode.IsUpper(ch) ||
		unicode.IsLower(ch) ||
		ch == '.' ||
		ch == '~' ||
		ch == '^' ||
		ch == '_' ||
		ch == ':' ||
		ch == '+'
}

func isYumModuleVersionFixed(mod string) bool {
	return strings.Contains(mod, ":")
}
