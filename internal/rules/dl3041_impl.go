package rules

import (
	"strings"
	"unicode"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3041 checks for dnf/microdnf install without version pinning.
func DL3041() rule.Rule {
	return rule.NewSimpleRule(
		DL3041Meta.Code,
		DL3041Meta.Severity,
		DL3041Meta.Message,
		checkDL3041,
	)
}

func checkDL3041(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	// Check all dnf/microdnf packages
	packages := getDnfPackages(parsed)
	for _, pkg := range packages {
		if !isDnfPackageVersionFixed(pkg) {
			return false
		}
	}

	// Check all dnf/microdnf modules
	modules := getDnfModules(parsed)
	for _, mod := range modules {
		if !isDnfModuleVersionFixed(mod) {
			return false
		}
	}

	return true
}

func getDnfPackages(parsed *shell.ParsedShell) []string {
	var packages []string

	for _, cmd := range parsed.PresentCommands {
		// Skip module commands
		if isDnfModuleCmd(cmd) {
			continue
		}

		// Get packages from install filter
		packages = append(packages, dnfInstallFilter(cmd)...)
	}

	return packages
}

func getDnfModules(parsed *shell.ParsedShell) []string {
	var modules []string

	for _, cmd := range parsed.PresentCommands {
		// Only module commands
		if !isDnfModuleCmd(cmd) {
			continue
		}

		// Get modules from install filter
		modules = append(modules, dnfInstallFilter(cmd)...)
	}

	return modules
}

func isDnfModuleCmd(cmd shell.Command) bool {
	if cmd.Name != "dnf" && cmd.Name != "microdnf" {
		return false
	}

	return shell.CmdHasArgs(cmd.Name, []string{"module"}, cmd)
}

func dnfInstallFilter(cmd shell.Command) []string {
	if cmd.Name != "dnf" && cmd.Name != "microdnf" {
		return nil
	}

	// Must be install command
	if !shell.CmdHasArgs(cmd.Name, []string{"install"}, cmd) {
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

func isDnfPackageVersionFixed(pkg string) bool {
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
	return isDnfVersionLike(versionParts)
}

func isDnfVersionLike(parts []string) bool {
	if len(parts) == 0 {
		return false
	}

	allValid := true
	hasDigitStart := false

	for _, part := range parts {
		if !isDnfValidVersionPart(part) {
			allValid = false
			break
		}
		if len(part) > 0 && unicode.IsDigit(rune(part[0])) {
			hasDigitStart = true
		}
	}

	return allValid && hasDigitStart
}

func isDnfValidVersionPart(part string) bool {
	if len(part) == 0 {
		return false
	}

	for _, ch := range part {
		if !isDnfVersionChar(ch) {
			return false
		}
	}
	return true
}

func isDnfVersionChar(ch rune) bool {
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

func isDnfModuleVersionFixed(mod string) bool {
	return strings.Contains(mod, ":")
}
