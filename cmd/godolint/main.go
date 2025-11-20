// Package main implements the godolint CLI for linting Dockerfiles.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return errors.New("at least one argument required: path to Dockerfile(s)")
			}

			// Build rule set
			rules := sdk.AllRules()

			// Add shellcheck by default unless disabled or binary is missing
			if !cmd.Bool("without-shellcheck") {
				// Check if shellcheck binary exists
				if _, err := exec.LookPath("shellcheck"); err != nil {
					// Shellcheck binary not found, print warning
					log.Warn().Msg("shellcheck binary not found in PATH, shellcheck integration disabled")
				} else {
					// Shellcheck available, use it
					checker := shell.NewBinaryShellchecker()
					scRule := shell.NewShellcheckRule(checker)
					rules = append(rules, scRule)
				}
			}

			// Create processor with all rules (reuse for all files)
			processor := process.NewProcessor(rules).
				WithDisableIgnorePragmas(cmd.Bool("disable-ignore-pragma"))

			// Collect all failures from all files
			allFailures := []rule.CheckFailure{}

			// Process each file
			for _, dockerfilePath := range cmd.Args().Slice() {
				// Read Dockerfile
				dockerfileContent, err := os.ReadFile(dockerfilePath)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", dockerfilePath, err)
				}

				// Parse Dockerfile using buildkit parser
				p := parser.NewBuildkitParser()
				instructions, err := p.Parse(dockerfileContent)
				if err != nil {
					return fmt.Errorf("failed to parse %s: %w", dockerfilePath, err)
				}

				log.Debug().Str("file", dockerfilePath).Int("instructions", len(instructions)).Msg("Parsed Dockerfile")

				// Run rules
				failures := processor.Run(instructions)

				// Add file path to each failure
				for i := range failures {
					failures[i].File = dockerfilePath
				}

				// Collect failures
				allFailures = append(allFailures, failures...)
			}

			// Filter out ignored rules
			ignoredRules := cmd.StringSlice("ignore")
			if len(ignoredRules) > 0 {
				filtered := []rule.CheckFailure{}
				for _, failure := range allFailures {
					ignored := false
					for _, ignoreCode := range ignoredRules {
						if string(failure.Code) == ignoreCode {
							ignored = true

							break
						}
					}
					if !ignored {
						filtered = append(filtered, failure)
					}
				}
				allFailures = filtered
			}

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

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Error().Err(err).Msg("failed to run godolint")
		os.Exit(1)
	}
}
