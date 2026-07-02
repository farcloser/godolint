// Package main implements the godolint CLI for linting Dockerfiles.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/sdk"
)

// errUsage reports an invocation without any Dockerfile argument.
var errUsage = errors.New("at least one argument required: path to Dockerfile(s)")

// buildRules assembles the rule set, wiring in the shellcheck integration
// unless it is disabled or the binary is missing from PATH.
func buildRules(cmd *cli.Command) ([]rule.Rule, error) {
	rules := sdk.AllRules()

	if cmd.Bool("without-shellcheck") {
		return rules, nil
	}

	// Missing shellcheck degrades gracefully (matching hadolint): warn and
	// lint without the integration rather than failing the run.
	//nolint:nilerr // intentional, see above.
	if _, err := exec.LookPath("shellcheck"); err != nil {
		log.Warn().Msg("shellcheck binary not found in PATH, shellcheck integration disabled")

		return rules, nil
	}

	checker := shell.NewBinaryShellchecker()
	// Fail fast on an unreadable rcfile: shellcheck errors are non-fatal per
	// rule (matching hadolint), so a bad path would otherwise silently
	// disable every SC check.
	if rcfile := cmd.String("shellcheck-rcfile"); rcfile != "" {
		if _, err := os.Stat(rcfile); err != nil {
			return nil, fmt.Errorf("cannot read shellcheck rcfile: %w", err)
		}

		checker.RCFile = rcfile
	}

	return append(rules, shell.NewShellcheckRule(checker)), nil
}

// lintFiles runs the processor over each Dockerfile and returns the collected
// failures, each tagged with the file it came from.
func lintFiles(processor *process.Processor, paths []string) ([]rule.CheckFailure, error) {
	// Non-nil so an all-clean run still encodes as JSON [] rather than null.
	allFailures := []rule.CheckFailure{}

	for _, dockerfilePath := range paths {
		//nolint:gosec // G304: reading user-supplied Dockerfile paths is this tool's purpose.
		dockerfileContent, err := os.ReadFile(dockerfilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", dockerfilePath, err)
		}

		instructions, err := parser.NewBuildkitParser().Parse(dockerfileContent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", dockerfilePath, err)
		}

		log.Debug().Str("file", dockerfilePath).Int("instructions", len(instructions)).Msg("Parsed Dockerfile")

		failures := processor.Run(instructions)
		for i := range failures {
			failures[i].File = dockerfilePath
		}

		allFailures = append(allFailures, failures...)
	}

	return allFailures, nil
}

// dropIgnored filters out failures whose rule code was --ignore'd.
func dropIgnored(failures []rule.CheckFailure, ignoredRules []string) []rule.CheckFailure {
	if len(ignoredRules) == 0 {
		return failures
	}

	filtered := []rule.CheckFailure{}

	for _, failure := range failures {
		if !slices.Contains(ignoredRules, string(failure.Code)) {
			filtered = append(filtered, failure)
		}
	}

	return filtered
}

func configureLogger(ctx context.Context, level ...zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger.WithContext(ctx)

	if len(level) > 0 {
		// Explicit level provided
		zerolog.SetGlobalLevel(level[0])
	} else {
		// Read from LOG_LEVEL environment variable
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}

		parsedLevel, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			// Invalid level, default to info
			parsedLevel = zerolog.InfoLevel

			log.Warn().Str("LOG_LEVEL", logLevel).Msg("Invalid log level, defaulting to info")
		}

		zerolog.SetGlobalLevel(parsedLevel)
	}
}

func main() {
	ctx := context.Background()
	configureLogger(ctx)

	cmd := &cli.Command{
		Name:  "godolint",
		Usage: "Dockerfile linter",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "disable-ignore-pragma",
				Usage: "Disable inline ignore pragmas `# hadolint ignore=DLxxxx`",
			},
			&cli.BoolFlag{
				Name:  "without-shellcheck",
				Usage: "Disable shellcheck integration for RUN instruction validation",
			},
			&cli.StringSliceFlag{
				Name:  "ignore",
				Usage: "Rule code to ignore (can be specified multiple times, e.g., --ignore DL3006 --ignore SC2050)",
			},
			&cli.StringFlag{
				Name:  "shellcheck-rcfile",
				Usage: "Shellcheckrc `FILE` forwarded to shellcheck (--rcfile) when validating RUN instructions (requires shellcheck >= 0.10.0)",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return errUsage
			}

			rules, err := buildRules(cmd)
			if err != nil {
				return err
			}

			// Create processor with all rules (reuse for all files)
			processor := process.NewProcessor(rules).
				WithDisableIgnorePragmas(cmd.Bool("disable-ignore-pragma"))

			allFailures, err := lintFiles(processor, cmd.Args().Slice())
			if err != nil {
				return err
			}

			allFailures = dropIgnored(allFailures, cmd.StringSlice("ignore"))

			// Output failures as JSON
			if err := json.NewEncoder(os.Stdout).Encode(allFailures); err != nil {
				return fmt.Errorf("failed to encode failures: %w", err)
			}

			// Exit with code 1 if any failures found
			if len(allFailures) > 0 {
				os.Exit(1)
			}

			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Error().Err(err).Msg("failed to run godolint")
		os.Exit(1)
	}
}
