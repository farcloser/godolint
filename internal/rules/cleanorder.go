package rules

import "github.com/forkcloser/godolint/internal/shell"

// cleanFollowsInstall is hadolint's ordering test for the package-manager
// cleanup rules: the first command install matches comes before the first
// command clean matches. False when either is absent, so a cache cleaned
// before the install that fills it counts as not cleaned.
func cleanFollowsInstall(cmds []shell.Command, install, clean func(shell.Command) bool) bool {
	installAt, cleanAt := firstCommand(cmds, install), firstCommand(cmds, clean)

	return installAt >= 0 && cleanAt >= 0 && installAt < cleanAt
}

// firstCommand is the index of the first command match holds for, or -1.
func firstCommand(cmds []shell.Command, match func(shell.Command) bool) int {
	for idx, cmd := range cmds {
		if match(cmd) {
			return idx
		}
	}

	return -1
}
