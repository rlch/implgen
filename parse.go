package main

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"io"
	"io/fs"
	"log/slog"
	"path"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	tsgo "github.com/smacker/go-tree-sitter/golang"
)

type (
	Repository struct {
		Package     string
		PackagePath string
		Filename    string
		Ident       string
		Methods     []*Method
		Imports     []Import
	}
	RepositoryImpl struct {
		Repository
		IsNew           bool
		ImplPackage     string
		ImplPackagePath string
		ImplFilename    string
		ImplMethods     []string
	}
	Import struct {
		Name string
		Path string
	}
	Method struct {
		Ident   string
		Params  Params
		Returns Params
	}
	Params []*Param
	Param  struct {
		Ident string
		Type  string
	}
)

var (
	ErrNoPackage = errors.New("no package name found")

	tsparser *sitter.Parser
	lang     = tsgo.GetLanguage()
)

func init() {
	tsparser = sitter.NewParser()
	tsparser.SetLanguage(lang)
}

func parseRepositoriesForPackage(
	ctx context.Context,
	fsys fs.FS,
	packagePath string,
	packageFiles []string,
) (repos []*Repository, err error) {
	repos = []*Repository{}
	for _, filename := range packageFiles {
		fullPath := path.Join(packagePath, filename)
		file, err := fsys.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fullPath, err)
		}
		defer file.Close()
		var src []byte
		src, err = io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
		}
		tree, err := tsparser.ParseCtx(ctx, nil, src)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", fullPath, err)
		}
		packageRepos, err := parseRepositories(src, tree)
		if err != nil {
			return nil, fmt.Errorf("failed to extract repositories from file %s: %w", fullPath, err)
		}
		for _, repo := range packageRepos {
			repo.Filename = filename
			repo.PackagePath = packagePath
		}
		repos = append(repos, packageRepos...)
	}
	return repos, nil
}

func parseRepositories(src []byte, tree *sitter.Tree) (repos []*Repository, err error) {
	dstFile, err := parser.ParseFile(
		fset,
		"",
		src,
		parser.ImportsOnly,
	)
	if err != nil {
		return nil, err
	}
	imports := make([]Import, len(dstFile.Imports))
	for i, imp := range dstFile.Imports {
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		}
		path, _ := strconv.Unquote(imp.Path.Value)
		imports[i] = Import{
			Name: name,
			Path: path,
		}
	}
	defer func() {
		for _, repo := range repos {
			repo.Imports = imports
		}
	}()

	const (
		PKG_CAPTURE = iota
		CLASS_NAME_CAPTURE
		METHOD_NAME_CAPTURE
		PARAMS_CAPTURE
		RESULT_CAPTURE
	)
	query, err := sitter.NewQuery([]byte(`
(package_clause (package_identifier) @pkg) 

(type_spec
  name: (type_identifier) @class_name (#match? @class_name "Repository$")
  type: 
   (interface_type
     (method_elem
       name: (field_identifier) @method_name
       parameters: (parameter_list) @params
       result: [
        (parameter_list)
        (type_identifier)
       ]? @result)?))
    `), lang)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	qc := sitter.NewQueryCursor()
	qc.Exec(query, tree.RootNode())

	// Get package name
	m, ok := qc.NextMatch()
	if !ok {
		return nil, nil
	}
	m = qc.FilterPredicates(m, src)
	if len(m.Captures) != 1 {
		return nil, ErrNoPackage
	}

	pkg := m.Captures[0].Node.Content(src)
	defer func() {
		for _, repo := range repos {
			repo.Package = pkg
		}
	}()
	m, ok = qc.NextMatch()
	if !ok {
		return nil, nil
	}
	var curIdx, methodIdx int
	repos = append(repos, &Repository{})
	for {
		m = qc.FilterPredicates(m, src)
		for _, c := range m.Captures {
			repo := repos[curIdx]
			switch c.Index {
			case CLASS_NAME_CAPTURE:
				name := c.Node.Content(src)
				if repo.Ident == name {
					continue
				} else if repo.Ident != "" {
					// We have a new repository, so append the previous one to the list.
					methodIdx = 0
					curIdx++
					repos = append(repos, &Repository{Ident: name})
				} else if repo.Ident == "" {
					repo.Ident = name
				}
			case METHOD_NAME_CAPTURE:
				var curMethod *Method
				methodName := c.Node.Content(src)
				if len(repo.Methods) == 0 {
					curMethod = &Method{Ident: methodName}
					repo.Methods = append(repo.Methods, curMethod)
				} else {
					curMethod = repo.Methods[methodIdx]
				}
				if curMethod.Ident != methodName {
					methodIdx++
					repo.Methods = append(repo.Methods, &Method{Ident: methodName})
				}

			case PARAMS_CAPTURE:
				repo.Methods[methodIdx].Params = parseParams(c.Node.Content(src))
			case RESULT_CAPTURE:
				repo.Methods[methodIdx].Returns = parseParams(c.Node.Content(src))
			default:
				slog.Error(
					"unhandled",
					slog.Int("index", int(c.Index)),
					slog.String("src", c.Node.Content(src)),
				)
			}
		}
		m, ok = qc.NextMatch()
		if !ok {
			break
		}
	}
	return
}

func parseParams(src string) []*Param {
	if src == "" {
		return nil
	}
	if src[0] == '(' {
		src = src[1 : len(src)-1]
	}
	parts := strings.Split(strings.TrimSpace(src), ",")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		return nil
	}
	named := false
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if named || len(strings.Split(part, " ")) > 1 {
			named = true
		}
	}
	if named {
		return parseNamedParams(parts)
	}
	params := make([]*Param, len(parts))
	for i, part := range parts {
		params[i] = &Param{Type: strings.TrimSpace(part)}
	}
	return params
}

func parseNamedParams(parts []string) []*Param {
	params := make([]*Param, len(parts))
	untypedFrom := -1
	for i, part := range parts {
		part = strings.TrimSpace(part)
		ident, typ, hasType := strings.Cut(part, " ")
		ident, typ = strings.TrimSpace(ident), strings.TrimSpace(typ)
		param := &Param{Ident: ident}
		params[i] = param
		if hasType {
			param.Type = typ
			if untypedFrom != -1 {
				for j := untypedFrom; j < i; j++ {
					params[j].Type = typ
				}
				untypedFrom = -1
			}
		} else if untypedFrom == -1 {
			untypedFrom = i
		}
	}
	return params
}

func parseRepositoryImpls(
	ctx context.Context,
	fsys fs.FS,
	implPackagePath string,
	repos []*Repository,
) ([]*RepositoryImpl, error) {
	if len(repos) == 0 {
		return nil, nil
	}

	defaultImplFilename := func(repo *RepositoryImpl) string {
		return strings.ToLower(repo.Name()) + "_impl.go"
	}
	implPackageName := repos[0].Package + "impl"
	impls := make([]*RepositoryImpl, len(repos))
	for i, repo := range repos {
		impls[i] = &RepositoryImpl{
			Repository: *repo,
		}
	}
	entries, err := fs.ReadDir(fsys, implPackagePath)
	if err != nil {
		err := err.(*fs.PathError)
		if errors.Is(err.Err, fs.ErrNotExist) {
			for _, repo := range impls {
				repo.ImplPackage = implPackageName
				repo.ImplPackagePath = implPackagePath
				repo.ImplFilename = defaultImplFilename(repo)
				repo.IsNew = true
			}
			return impls, nil
		}
		return nil, err
	}

	implDeclsToFileMap := make(map[string]string)
	repositoryToMethodMap := make(map[string][]string)
	for _, d := range entries {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			continue
		}
		fi, err := d.Info()
		if err != nil {
			return nil, err
		}
		filename := fi.Name()
		if strings.HasSuffix(filename, "_test.go") || !strings.HasSuffix(filename, ".go") {
			continue
		}
		path := path.Join(implPackagePath, filename)
		file, err := fsys.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		src, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}
		pkg, implDecls, methods, err := parseRepositoryImplFile(ctx, src)
		if err != nil {
			return nil, err
		}
		implPackageName = pkg
		for _, decl := range implDecls {
			implDeclsToFileMap[decl] = filename
		}
		for rep, methods := range methods {
			repositoryToMethodMap[rep] = append(repositoryToMethodMap[rep], methods...)
		}
	}

	for _, repo := range impls {
		implName := repo.ImplName()
		repo.ImplPackage = implPackageName
		repo.ImplPackagePath = implPackagePath
		if filename, ok := implDeclsToFileMap[implName]; ok {
			repo.ImplFilename = filename
			repo.ImplMethods = repositoryToMethodMap[implName]
		} else {
			repo.ImplFilename = defaultImplFilename(repo)
			repo.IsNew = true
		}
	}
	return impls, nil
}

func parseRepositoryImplFile(ctx context.Context, src []byte) (
	packageName string,
	repImpls []string,
	methods map[string][]string,
	err error,
) {
	const (
		PKG_CAPTURE = iota
		IMPL_NAME_CAPTURE
		IMPL_REC_CAPTURE
		IMPL_FIELD_CAPTURE
	)
	repImpls = []string{}
	methods = make(map[string][]string)
	query, err := sitter.NewQuery([]byte(`
  (
    (package_clause (package_identifier) @pkg)
    (type_declaration 
      (type_spec
          name: (type_identifier) @impl_name (#match? @impl_name "Impl$")))?
    (method_declaration
        receiver: (parameter_list
          (parameter_declaration
            type: (_) @impl_rec (#match? @impl_rec "Impl$")))
        name: (field_identifier) @impl_field)?
  )
`), lang)
	if err != nil {
		return "", nil, nil, err
	}

	tree, err := tsparser.ParseCtx(ctx, nil, src)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to parse file: %w", err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(query, tree.RootNode())

	// Get package name
	m, ok := qc.NextMatch()
	if !ok {
		return "", nil, nil, nil
	}
	repImplMap := map[string]bool{}
	var curRec string
	for {
		m = qc.FilterPredicates(m, src)
		for _, c := range m.Captures {
			switch c.Index {
			case PKG_CAPTURE:
				packageName = c.Node.Content(src)
			case IMPL_NAME_CAPTURE:
				implName := c.Node.Content(src)
				repImplMap[implName] = true
			case IMPL_REC_CAPTURE:
				rec := c.Node.Content(src)
				if rec == "" {
					continue
				} else if rec[0] == '*' {
					rec = rec[1:]
				}
				curRec = rec
			case IMPL_FIELD_CAPTURE:
				if curRec == "" {
					panic("receiver not found")
				}
				method := c.Node.Content(src)
				found := false
				for _, curMethod := range methods[curRec] {
					if curMethod == method {
						found = true
						break
					}
				}
				if !found {
					methods[curRec] = append(methods[curRec], method)
				}
			default:
				slog.Error(
					"unhandled",
					slog.Int("index", int(c.Index)),
					slog.String("src", c.Node.Content(src)),
				)
			}
		}
		m, ok = qc.NextMatch()
		if !ok {
			break
		}
	}
	for impl := range repImplMap {
		repImpls = append(repImpls, impl)
	}
	return
}
