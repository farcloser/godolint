package rules

import "strings"

// hasCacheOrTmpfsMount checks if RUN flags contain a cache or tmpfs mount
// for the specified path.
func hasCacheOrTmpfsMount(flags []string, path string) bool {
	for _, flag := range flags {
		if !strings.HasPrefix(flag, "--mount=") {
			continue
		}

		mountSpec := strings.TrimPrefix(flag, "--mount=")
		parts := strings.Split(mountSpec, ",")

		mountType := ""
		target := ""

		for _, part := range parts {
			keyVal := strings.SplitN(part, "=", 2)
			if len(keyVal) != 2 {
				continue
			}

			key := strings.TrimSpace(keyVal[0])
			value := strings.TrimSpace(keyVal[1])

			switch key {
			case "type":
				mountType = value
			case "target":
				target = value
			}
		}

		// Check if it's a cache or tmpfs mount
		if mountType != "cache" && mountType != "tmpfs" {
			continue
		}

		// Check if target matches the path
		if target == path || strings.HasPrefix(target, path+"/") {
			return true
		}
	}

	return false
}
