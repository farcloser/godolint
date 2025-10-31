package rules

import (
	"strconv"
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// DL3011 creates a rule for checking port ranges are valid.
func DL3011() rule.Rule {
	return rule.NewSimpleRule(
		DL3011Meta.Code,
		DL3011Meta.Severity,
		DL3011Meta.Message,
		checkDL3011,
	)
}

func checkDL3011(instruction syntax.Instruction) bool {
	expose, ok := instruction.(*syntax.Expose)
	if !ok {
		return true
	}

	// Check all ports are in valid range (0-65535)
	for _, portSpec := range expose.Ports {
		if !isValidPortSpec(portSpec) {
			return false
		}
	}

	return true
}

// Port spec can be: "80", "80/tcp", "8000-9000", "8000-9000/tcp", "${VAR}", "8000-${VAR}".
func isValidPortSpec(portSpec string) bool {
	// Remove protocol suffix if present (/tcp, /udp, /sctp)
	parts := strings.Split(portSpec, "/")
	portPart := parts[0]

	// Variables are always OK
	if strings.Contains(portPart, "$") {
		return true
	}

	// Check if it's a range
	if strings.Contains(portPart, "-") {
		rangeParts := strings.Split(portPart, "-")
		if len(rangeParts) != 2 {
			return false
		}

		low, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
		high, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))

		if err1 != nil || err2 != nil {
			return false
		}

		return low >= 0 && low <= 65535 && high >= 0 && high <= 65535
	}

	// Single port
	port, err := strconv.Atoi(strings.TrimSpace(portPart))
	if err != nil {
		return false
	}

	return port >= 0 && port <= 65535
}
