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

// crawlAPI returns the list of files that should be analysed for API definitions.
// The returned map is a mapping of directories to the file names in that directory.
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

var cachedModule string

// getModule recursively searches (upwards) for a go.mod file and returns the module path.
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

// loadLocalPackage returns the import path and alias for a local package under
// the context of a (potentially already dependent) astFile.
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
	base := path.Base(packagePath)
	importAlias = strings.ReplaceAll(base, "_", "")
	if base == importAlias {
		importAlias = ""
	}
	return
}
