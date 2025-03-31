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

	fset = token.NewFileSet()
)

var (
	fRoot   string
	fApi    string
	fImpl   string
	fFocus  []string
	fUseDig bool

	verbose bool
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error(
			"Failed to run implgen",
			slog.Any("error", err),
		)
	}
}
