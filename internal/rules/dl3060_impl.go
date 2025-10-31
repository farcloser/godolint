package rules

import (
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3060 checks for yarn cache clean after yarn install.
func DL3060() rule.Rule {
	return rule.NewSimpleRule(
		DL3060Meta.Code,
		DL3060Meta.Severity,
		DL3060Meta.Message,
		checkDL3060,
	)
}

func checkDL3060(instruction syntax.Instruction) bool {
	run, ok := instruction.(*syntax.Run)
	if !ok {
		return true
	}

	// Check if cache/tmpfs mount is present for yarn cache
	if hasCacheOrTmpfsMount(run.Flags, ".cache/yarn") ||
		hasCacheOrTmpfsMount(run.Flags, "/root/.cache/yarn") {
		return true
	}

	parsed, err := shell.ParseShell(run.Command)
	if err != nil {
		return true
	}

	hasYarnInstall := false
	hasYarnCacheClean := false

	for _, cmd := range parsed.PresentCommands {
		if shell.CmdHasArgs("yarn", []string{"install"}, cmd) {
			hasYarnInstall = true
		}
		if shell.CmdHasArgs("yarn", []string{"cache", "clean"}, cmd) {
			hasYarnCacheClean = true
		}
	}

	// If no yarn install, pass
	if !hasYarnInstall {
		return true
	}

	// If has yarn install, must have cache clean
	return hasYarnCacheClean
}
