package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/urfave/cli/v3"
)

// generate orchestrates the main generation process for creating implementation files.
//
// It performs the following steps:
//  1. Crawls the API directory to find Go files containing interfaces
//  2. Applies focus filters using glob patterns if specified
//  3. Parses Repository interfaces from the discovered files
//  4. Computes implementation package paths and generates implementation files
//  5. Creates or updates the dependency injection stub file (repositories.go)
//
// The function preserves existing implementations and only generates missing methods.
func generate(ctx context.Context, cmd *cli.Command) error {
	fsys := os.DirFS(fRoot)
	slog.Debug(
		"Crawling API directory",
		slog.String("root", fRoot),
		slog.String("api_root", fApi),
		slog.String("impl_root", fImpl),
	)
	apiFiles, err := crawlAPI(fsys, fApi)
	if err != nil {
		return fmt.Errorf("failed to walk API directory: %w", err)
	}
	allRepImpls := []*RepositoryImpl{}
	focusGlobs := make([]glob.Glob, len(fFocus))
	for i, focus := range fFocus {
		glob, err := glob.Compile(focus)
		if err != nil {
			return fmt.Errorf("failed to compile glob %s: %w", focus, err)
		}
		focusGlobs[i] = glob
	}

	for apiPackagePath, packageFiles := range apiFiles {
		match := len(fFocus) == 0
		for _, glob := range focusGlobs {
			relToAPI, err := filepath.Rel(fApi, apiPackagePath)
			if err != nil {
				return fmt.Errorf("failed to compute relative path to API: %w", err)
			}
			if glob.Match(relToAPI) {
				match = true
				break
			}
		}
		if !match {
			slog.Debug(
				"Skipping package",
				slog.String("api_path", apiPackagePath),
			)
			continue
		}
		repos, err := parseRepositoriesForPackage(
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
			fApi,
			fImpl,
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
			fsys,
			implPackagePath,
			repos,
		)
		if err != nil {
			return fmt.Errorf("failed to parse repository implementations: %w", err)
		}
		allRepImpls = append(allRepImpls, repImpls...)
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
	}
	if len(fFocus) == 0 {
		stubSrc, err := generateRepositoryStubFile(fsys, fImpl, allRepImpls...)
		if err != nil {
			return fmt.Errorf("failed to generate repository stub file: %w", err)
		}
		if err := os.WriteFile(
			path.Join(fImpl, "repositories.go"),
			[]byte(stubSrc),
			0644,
		); err != nil {
			return fmt.Errorf("failed to write repository stub file: %w", err)
		}
		slog.Debug("Generated repository stub file")
	} else {
		slog.Debug("Focus provided, skipping stub generation")
	}
	return nil
}

// groupByPackage groups repository implementations by their API package name.
// This is used to organize repositories for mock generation directives.
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

// groupByImplFilename groups repository implementations by their implementation filename.
// This ensures that repositories that should be in the same file are processed together.
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
