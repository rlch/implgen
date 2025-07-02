// Package main provides the implgen CLI tool for generating Go interface implementations.
//
// implgen automatically creates implementation boilerplate for Go interfaces following
// the Repository pattern. It parses interface definitions from an API directory and
// generates corresponding implementation files with proper dependency injection,
// observability, and error handling.
package main

import (
	"context"
	"go/token"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
)

var (
	// cmd defines the root CLI command and its subcommands, flags, and behavior.
	cmd = &cli.Command{
		Name:           "implgen",
		Description:    "Code generator for API implementations.",
		DefaultCommand: "generate",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "Enable verbose logging",
				Aliases:     []string{"v"},
				Destination: &verbose,
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "generate",
				Description: "Generate API implementations",
				Action:      generate,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "root",
						Usage:       "Root directory to generate the api/impl tree from.",
						Value:       ".",
						Destination: &fRoot,
					},
					&cli.StringFlag{
						Name:        "api",
						Usage:       "Directory to API definitions, relative to root.",
						Value:       "api",
						Destination: &fApi,
					},
					&cli.StringFlag{
						Name:        "impl",
						Usage:       "Directory to implementation files, relative to root.",
						Value:       "internal",
						Destination: &fImpl,
					},
					&cli.StringSliceFlag{
						Name:        "focus",
						Usage:       "Focus generating specific packages relative to API using a glob. Does not generate stub if provided.",
						Destination: &fFocus,
					},
					&cli.BoolFlag{
						Name:        "dig",
						Usage:       "Use dig for dependency injection instead of fx.",
						Destination: &fUseDig,
					},
				},
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			logOpts := &tint.Options{
				TimeFormat: time.Kitchen,
			}
			if verbose {
				logOpts.Level = slog.LevelDebug
				logOpts.AddSource = true
			}
			logger := slog.New(
				tint.NewHandler(os.Stdout, logOpts),
			)
			slog.SetDefault(logger)
			return ctx, nil
		},
	}

	// fset is a global file set used for parsing Go source files.
	fset = token.NewFileSet()
)

var (
	// CLI flag variables that store user-provided configuration.
	fRoot   string   // Root directory for the project
	fApi    string   // API directory relative to root
	fImpl   string   // Implementation directory relative to root
	fFocus  []string // Glob patterns to focus on specific packages
	fUseDig bool     // Whether to use dig instead of fx for dependency injection

	verbose bool // Whether to enable verbose logging
)

// main is the entry point for the implgen CLI application.
// It sets up the CLI command and executes it with the provided arguments.
func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error(
			"Failed to run implgen",
			slog.Any("error", err),
		)
	}
}
