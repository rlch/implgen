// Package main contains file system utilities for directory traversal and path management.
//
// This file provides functions for crawling directory trees to find Go source files,
// computing implementation paths from API paths, detecting Go module information,
// and managing import paths for generated code.
package main

import (
	"errors"
	"go/ast"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
)

// crawlAPI traverses the API directory tree to find Go source files.
//
// It returns a map where keys are directory paths and values are slices of Go filenames
// within those directories. Test files (*_test.go) are excluded from the results.
// The function only considers files with the .go extension.
func crawlAPI(
	fsys fs.FS,
	apiDir string,
) (files map[string][]string, err error) {
	apiDir = path.Clean(apiDir)
	files = map[string][]string{}
	err = fs.WalkDir(fsys, apiDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d == nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		filename := path.Base(p)
		dir := path.Dir(p)
		files[dir] = append(files[dir], filename)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

// computeImplPackagePath calculates the implementation package path from an API package path.
//
// It takes the relative path of the API package from the API root and mirrors that
// structure under the implementation root. For example:
//   - apiRoot: "api", implRoot: "internal", apiPackagePath: "api/user"
//   - Returns: "internal/user"
func computeImplPackagePath(apiRoot, implRoot, apiPackagePath string) (string, error) {
	apiRoot = path.Clean(apiRoot)
	implRoot = path.Clean(implRoot)
	apiPackagePath = path.Clean(apiPackagePath)

	implPackagePath := implRoot
	apiPackageParts := strings.Split(apiPackagePath, "/")
	apiRootParts := strings.Split(apiRoot, "/")
	for i, part := range apiPackageParts {
		if i < len(apiRootParts) && apiRoot != "." {
			rootPart := apiRootParts[i]
			if rootPart != part {
				return "", errors.New("apiPackagePath is not nested under apiRoot")
			}
			continue
		}
		implPackagePath = path.Join(implPackagePath, part)
	}
	return implPackagePath, nil
}

// cachedModule stores the detected module path to avoid repeated filesystem operations.
var cachedModule string

// getModule recursively searches upward for a go.mod file and returns the module path.
//
// It starts from the given root directory and walks up the directory tree until
// it finds a go.mod file. The module path is cached for subsequent calls to
// improve performance. Returns an error if no go.mod file is found.
func getModule(fsys fs.FS, root string) (module string, err error) {
	if cachedModule != "" {
		return cachedModule, nil
	}
	defer func() {
		cachedModule = module
	}()
	curDir := path.Clean(root)
	file, err := fsys.Open(path.Join(curDir, "go.mod"))
	if errors.Is(err, fs.ErrNotExist) {
		if curDir == "." || curDir == "/" {
			return "", errors.New("could not find a go.mod in current or parent directory")
		}
		return getModule(fsys, path.Dir(curDir))
	} else if err != nil {
		return "", err
	}
	gomodBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	gomod, err := modfile.ParseLax("go.mod", gomodBytes, nil)
	if err != nil {
		return "", err
	}
	module = gomod.Module.Mod.Path
	return
}

// loadLocalPackage computes the full import path and alias for a local package.
//
// It combines the module path with the package path to create a full import path.
// If an AST file is provided, it checks for existing import aliases in that file.
// This is used when generating import statements in implementation files.
func loadLocalPackage(
	fsys fs.FS,
	astFile *ast.File,
	packagePath string,
) (
	importPath, importAlias string,
	err error,
) {
	module, err := getModule(fsys, packagePath)
	if err != nil {
		return "", "", err
	}
	importPath = path.Join(module, packagePath)

	// Check if package has an existing alias or use default alias.
	if astFile != nil {
		for _, imp := range astFile.Imports {
			var path string
			path, err = strconv.Unquote(imp.Path.Value)
			if err != nil || path != importPath {
				continue
			}
			if imp.Name != nil {
				importAlias = imp.Name.Name
			}
			return
		}
	}
	return
}
