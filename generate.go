package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	slog.Info(
		"Starting generation",
		slog.String("root", fRoot),
		slog.String("api_root", fApi),
		slog.String("impl_root", fImpl),
		slog.String("suffix", fSuffix),
	)
	fsys := os.DirFS(fRoot)
	slog.Debug(
		"Crawling API directory",
		slog.String("root", fRoot),
		slog.String("api_root", fApi),
		slog.String("impl_root", fImpl),
		slog.String("suffix", fSuffix),
	)
	apiFiles, err := crawlAPI(fsys, fApi)
	if err != nil {
		return fmt.Errorf("failed to walk API directory: %w", err)
	}
	slog.Info(
		"Found API files",
		slog.Int("package_count", len(apiFiles)),
	)
	allContractImpls := []*ContractImpl{}
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
		contracts, err := parseContractsForPackage(
			fsys,
			apiPackagePath,
			packageFiles,
			fSuffix,
		)
		if err != nil {
			return fmt.Errorf("failed to parse contracts in %s: %w", apiPackagePath, err)
		}
		if len(contracts) == 0 {
			slog.Info(
				"No contracts found in package",
				slog.String("api_path", apiPackagePath),
				slog.String("suffix", fSuffix),
			)
			continue
		}
		slog.Info(
			"Parsed contracts",
			slog.String("api_path", apiPackagePath),
			slog.Int("count", len(contracts)),
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
		contractImpls, err := parseContractImpls(
			fsys,
			implPackagePath,
			contracts,
		)
		if err != nil {
			return fmt.Errorf("failed to parse contract implementations: %w", err)
		}
		allContractImpls = append(allContractImpls, contractImpls...)
		for filename, impls := range groupByImplFilename(contractImpls) {
			implPath := path.Join(implPackagePath, filename)
			fullImplPath := path.Join(fRoot, implPath)
			_, statErr := os.Stat(fullImplPath)
			exists := statErr == nil
			slog.Info(
				"Generating implementation file",
				slog.String("filename", filename),
				slog.String("impl_path", implPath),
				slog.Bool("exists", exists),
			)
			data, err := generateContractImplsForFile(fsys, implPath, impls)
			if err != nil {
				return fmt.Errorf("failed to generate implementation file: %w", err)
			}
			if data == "" {
				continue
			}
			if err := os.MkdirAll(
				path.Dir(fullImplPath),
				0755,
			); err != nil {
				return fmt.Errorf("failed to create directory for implementation file at %s: %w", fullImplPath, err)
			}
			if err := os.WriteFile(
				fullImplPath,
				[]byte(data),
				0644,
			); err != nil {
				return fmt.Errorf("failed to write implementation file at %s: %w", fullImplPath, err)
			}

			var nNewImpls, nNewMethods int
			for _, impl := range impls {
				if impl.IsNew {
					nNewImpls++
				}
				nNewMethods += len(impl.NewMethods())
			}
			if nNewImpls == 0 && nNewMethods == 0 {
				slog.Info(
					"No new methods or implementations to generate",
					slog.String("impl_path", implPath),
				)
				continue
			}
			var logMsg string
			if exists {
				logMsg = "Updated implementation file"
			} else {
				logMsg = "Created implementation file"
			}
			slog.Info(
				logMsg,
				slog.String("api_path", apiPackagePath),
				slog.String("impl_path", implPath),
				slog.Int("new_implementations", nNewImpls),
				slog.Int("new_methods", nNewMethods),
			)
		}
	}
	if len(fFocus) == 0 {
		stubSrc, err := generateContractStubFile(fsys, fImpl, fSuffix, allContractImpls...)
		if err != nil {
			return fmt.Errorf("failed to generate %s stub file: %w", strings.ToLower(fSuffix), err)
		}
		stubFileName := strings.ToLower(fSuffix) + ".go"
		stubFullPath := path.Join(fRoot, fImpl, stubFileName)
		if err := os.WriteFile(
			stubFullPath,
			[]byte(stubSrc),
			0644,
		); err != nil {
			return fmt.Errorf("failed to write %s stub file: %w", fSuffix, err)
		}
		slog.Info("Generated stub file", slog.String("path", path.Join(fImpl, stubFileName)))
	} else {
		slog.Debug("Focus provided, skipping stub generation")
	}
	return nil
}

// groupByPackage groups contract implementations by their API package name.
// This is used to organize contracts for mock generation directives.
func groupByPackage(contracts []*ContractImpl) map[string][]*ContractImpl {
	grouped := make(map[string][]*ContractImpl)
	for _, contract := range contracts {
		grouped[contract.Package] = append(
			grouped[contract.Package],
			contract,
		)
	}
	return grouped
}

// groupByImplFilename groups contract implementations by their implementation filename.
// This ensures that contracts that should be in the same file are processed together.
func groupByImplFilename(contracts []*ContractImpl) map[string][]*ContractImpl {
	grouped := make(map[string][]*ContractImpl)
	for _, contract := range contracts {
		grouped[contract.ImplFilename] = append(
			grouped[contract.ImplFilename],
			contract,
		)
	}
	return grouped
}
