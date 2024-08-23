package main

import (
	"context"
	"fmt"
	"go/token"
	"log/slog"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"
)

var (
	cli struct {
		Root    string `type:"path" help:"Root directory to generate the api/impl tree from." default:"."`
		API     string `type:"string" help:"Directory to API definitions, relative to root." default:"api"`
		Impl    string `type:"string" help:"Directory to implementation files, relative to root." default:"internal"`
		Verbose bool   `help:"Enable verbose logging." short:"v"`
	}
	fset = token.NewFileSet()
)

func main() {
	kong.Parse(
		&cli,
		kong.Name("implgen"),
		kong.Description("Code generator for API implementations."),
	)
	logOpts := &tint.Options{
		TimeFormat: time.Kitchen,
	}
	if cli.Verbose {
		logOpts.Level = slog.LevelDebug
		logOpts.AddSource = true
	}
	logger := slog.New(
		tint.NewHandler(os.Stdout, logOpts),
	)
	slog.SetDefault(logger)
	if err := run(); err != nil {
		logger.Error(
			"Failed to run implgen",
			slog.Any("error", err),
		)
		if cli.Verbose {
			debug.PrintStack()
		}
	}
}

func run() error {
	ctx := context.Background()
	fsys := os.DirFS(cli.Root)
	slog.Debug(
		"Crawling API directory",
		slog.String("root", cli.Root),
		slog.String("api_root", cli.API),
		slog.String("impl_root", cli.Impl),
	)
	apiFiles, err := crawlAPI(fsys, cli.API)
	if err != nil {
		return fmt.Errorf("failed to walk API directory: %w", err)
	}
	for apiPackagePath, packageFiles := range apiFiles {
		repos, err := parseRepositoriesForPackage(
			ctx,
			fsys,
			apiPackagePath,
			packageFiles,
		)
		if err != nil {
			return fmt.Errorf("failed to parse repositories in %s: %w", apiPackagePath, err)
		}
		if len(repos) == 0 {
			continue
		}
		slog.Debug(
			"Parsed repositories",
			slog.String("api_path", apiPackagePath),
			slog.Int("count", len(repos)),
		)
		implPackagePath, err := computeImplPackagePath(
			cli.API,
			cli.Impl,
			apiPackagePath,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to compute implementation package path associated with API %s: %w",
				apiPackagePath,
				err,
			)
		}
		repImpls, err := parseRepositoryImpls(
			ctx,
			fsys,
			implPackagePath,
			repos,
		)
		if err != nil {
			return fmt.Errorf("failed to parse repository implementations: %w", err)
		}
		for filename, impls := range groupByImplFilename(repImpls) {
			implPath := path.Join(implPackagePath, filename)
			_, statErr := os.Stat(implPath)
			exists := statErr == nil
			data, err := generateRepositoryImplsForFile(fsys, implPath, impls)
			if err != nil {
				return fmt.Errorf("failed to generate implementation file: %w", err)
			}
			if data == "" {
				continue
			}
			if err := os.MkdirAll(
				path.Dir(implPath),
				0755,
			); err != nil {
				return fmt.Errorf("failed to create directory for implementation file at %s: %w", implPath, err)
			}
			if err := os.WriteFile(
				implPath,
				[]byte(data),
				0644,
			); err != nil {
				return fmt.Errorf("failed to write implementation file at %s: %w", implPath, err)
			}

			var nNewImpls, nNewMethods int
			for _, impl := range impls {
				if impl.IsNew {
					nNewImpls++
				}
				nNewMethods += len(impl.NewMethods())
			}
			if nNewImpls == 0 && nNewMethods == 0 {
				continue
			}
			var logMsg string
			if exists {
				logMsg = "Updated implementation file"
			} else {
				logMsg = "Created implementation file"
			}
			slog.Debug(
				logMsg,
				slog.String("api_path", apiPackagePath),
				slog.String("impl_path", implPath),
				slog.Int("new_implementations", nNewImpls),
				slog.Int("new_methods", nNewMethods),
			)
		}
		stubSrc, err := generateRepositoryStubFile(fsys, cli.Impl, repImpls...)
		if err != nil {
			return fmt.Errorf("failed to generate repository stub file: %w", err)
		}
		if err := os.WriteFile(
			path.Join(cli.Impl, "repositories.go"),
			[]byte(stubSrc),
			0644,
		); err != nil {
			return fmt.Errorf("failed to write repository stub file: %w", err)
		}
		slog.Debug("Generated repository stub file")
	}
	return nil
}

func groupByPackage(repositories []*RepositoryImpl) map[string][]*RepositoryImpl {
	grouped := make(map[string][]*RepositoryImpl)
	for _, repository := range repositories {
		grouped[repository.Package] = append(
			grouped[repository.Package],
			repository,
		)
	}
	return grouped
}

func groupByImplFilename(repositories []*RepositoryImpl) map[string][]*RepositoryImpl {
	grouped := make(map[string][]*RepositoryImpl)
	for _, repository := range repositories {
		grouped[repository.ImplFilename] = append(
			grouped[repository.ImplFilename],
			repository,
		)
	}
	return grouped
}
