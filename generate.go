package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

func (r RepositoryImpl) NewMethods() []*Method {
	methods := []*Method{}
	for _, method := range r.Methods {
		existing := false
		for _, existingMethod := range r.ImplMethods {
			if method.Ident == existingMethod {
				existing = true
				break
			}
		}
		if existing {
			continue
		}
		args := make(Params, len(method.Params))
		returns := make(Params, len(method.Returns))
		qualify := func(arg *Param) *Param {
			isLower := 'a' <= arg.Type[0] && arg.Type[0] <= 'z'
			if strings.Contains(arg.Type, ".") || isLower {
				return arg
			}
			return &Param{
				Ident: arg.Ident,
				Type:  r.Package + "." + arg.Type,
			}
		}
		for i, arg := range method.Params {
			arg := arg
			args[i] = qualify(arg)
		}
		for i, arg := range method.Returns {
			arg := arg
			returns[i] = qualify(arg)
		}
		if len(args) == 0 {
			args = nil
		}
		if len(returns) == 0 {
			returns = nil
		}
		methods = append(methods, &Method{
			Ident:   method.Ident,
			Params:  args,
			Returns: returns,
		})
	}
	return methods
}

func (p Params) HasCtx() bool {
	for _, param := range p {
		if param.Type == "context.Context" {
			return true
		}
	}
	return false
}

func (p Params) HasError() bool {
	for _, param := range p {
		if param.Type == "error" {
			return true
		}
	}
	return false
}

func (p Params) Named() bool {
	for _, param := range p {
		if param.Ident != "" {
			return true
		}
	}
	return false
}

func (p Params) Qualify() {
	hasCtx := p.HasCtx()
	hasErr := p.HasError()
	named := p.Named()
	if !hasCtx && !hasErr && !named {
		return
	}
	for _, param := range p {
		switch param.Type {
		case "context.Context":
			param.Ident = "ctx"
		case "error":
			param.Ident = "err"
		default:
			if param.Ident == "" {
				param.Ident = "_"
			}
		}
	}
}

func (p Params) ParamsSrc() (s string) {
	p.Qualify()
	n := len(p)
	for i, param := range p {
		if i > 0 {
			s += ", "
		}
		if i < n-1 && param.Type == p[i+1].Type {
			s += param.Ident
		} else if param.Ident != "" {
			s += param.Ident + " " + param.Type
		} else {
			s += param.Type
		}
	}
	return
}

func (p Params) ReturnsSrc() string {
	p.Qualify()
	src := p.ParamsSrc()
	if p.Named() {
		return "(" + src + ")"
	}
	return src
}

func (r Repository) QualifyString(s string) string {
	name := r.Name()
	if name == "Repository" {
		return s
	}
	return name + s
}

func (r Repository) Name() string {
	if len(r.Ident) > 10 && strings.HasSuffix(r.Ident, "Repository") {
		return r.Ident[:len(r.Ident)-10]
	}
	return r.Ident
}

func (r Repository) ImplName() string {
	if r.Ident == "" {
		return ""
	}
	name := r.Ident
	return strings.ToLower(string(name[0])) + name[1:] + "Impl"
}

func (r Repository) QualifiedName() string {
	if r.Package == "" {
		return r.Ident
	}
	return r.Package + "." + r.Ident
}

const generateMethodTemplate = `
  func (r *{{ .Repository.ImplName }}) {{ .Method.Ident }}({{ .Method.Params.ParamsSrc }}){{ pad .Method.Returns.ReturnsSrc }}{
  {{- if .Method.Params.HasCtx }}
    ctx, span := otel.GetTracerProvider().Tracer("{{ .Repository.Package }}").Start(ctx, "{{ .Repository.Name }}.{{ .Method.Ident }}")
    {{- if .Method.Returns.HasError }}
    defer func() {
      if err != nil {
        err = eris.Wrap(err, "{{ .Repository.QualifiedName }}.{{ .Method.Ident }}")
        span.SetStatus(codes.Error, "")
        span.RecordError(err)
      }
      span.End()
    }()
    {{- else }}
    defer span.End()
    {{- end }}
    _ = ctx
  {{- else }}
    {{- if .Method.Returns.HasError }}
    defer func() {
      if err != nil {
        err = eris.Wrap(err, "{{ .Repository.QualifiedName }}.{{ .Method.Ident }}")
      }
    }()
    {{- end }}
  {{- end }}
    panic("TODO: implement {{ .Repository.QualifiedName }}.{{ .Method.Ident }}")
  }
`

func generateMethodImpl(repository Repository, method Method) (string, error) {
	tmpl, err := template.
		New("generateMethodTemplate").
		Funcs(template.FuncMap{
			"pad": func(s string) string {
				if s == "" {
					return " "
				}
				return " " + s + " "
			},
		}).
		Parse(generateMethodTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Repository
		Method
	}{repository, method}); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

const generateRepositoryImplTemplate = `
type {{ .Repository.QualifyString "Dependencies" }} struct {
  fx.In
	// Add dependencies here
}

var {{ .Repository.QualifyString "Options" }} = fx.Options(
	fx.Provide(
		New{{ .Repository.Ident }},
	),
)

func New{{ .Repository.Ident }}(deps {{ .Repository.QualifyString "Dependencies" }}) {{ .Repository.Package }}.{{ .Repository.Ident }} {
	return &{{ .Repository.ImplName }}{
    {{ .Repository.QualifyString "Dependencies" }}: deps,
	}
}

type {{ .Repository.ImplName }} struct {
  {{ .Repository.QualifyString "Dependencies" }}
}
`

// generateRepositoryImpl generates the method and struct declarations for a single repository.
func generateRepositoryImpl(repository Repository) (string, error) {
	tmpl, err := template.
		New("generateRepositoryImplTemplate").
		Parse(generateRepositoryImplTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Repository
	}{repository}); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// generateRepositoryImplsForFile generates the repository implementations for a single file.
//
// All RepositoryImpl's are assumed to be for the same file as repositories[0].
func generateRepositoryImplsForFile(
	fsys fs.FS,
	filepath string,
	repositories []*RepositoryImpl,
) (_ string, err error) {
	if len(repositories) == 0 {
		return "", nil
	}
	var (
		originalSrc        []byte
		originalSrcScanner *bufio.Scanner
		src                bytes.Buffer
		astFile            *ast.File
	)
	file, err := fsys.Open(filepath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	// Write package declaration to src. If the file does not exist, write a new package declaration.
	if file != nil {
		defer file.Close()
		originalSrc, err = io.ReadAll(file)
		if err != nil {
			return "", err
		}
		astFile, err = parser.ParseFile(fset, "", originalSrc, parser.ImportsOnly)
		if err != nil {
			return "", err
		}
		originalSrcScanner = bufio.NewScanner(bytes.NewReader(originalSrc))
		for originalSrcScanner.Scan() {
			line := originalSrcScanner.Text()
			src.WriteString(line + "\n")
			if strings.HasPrefix(line, "package") {
				break
			}
		}
	} else {
		packageDecl := fmt.Sprintf(`
// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package %s
`, repositories[0].ImplPackage)
		src.WriteString(strings.TrimPrefix(packageDecl, "\n"))
	}

	// Add imports to src
	requiredImports, err := collectImports(
		fsys,
		astFile,
		true,
		false,
		repositories...,
	)
	if err != nil {
		return "", err
	}
	for _, imp := range requiredImports {
		src.WriteString("import ")
		if imp.Name != "" {
			src.WriteString(imp.Name + " ")
		}
		src.WriteString(strconv.Quote(imp.Path) + "\n")
	}

	// Add the rest of the original source code.
	if originalSrcScanner != nil {
		for originalSrcScanner.Scan() {
			src.WriteString(originalSrcScanner.Text() + "\n")
		}
	}

	// Append new repository declarations
	for _, repository := range repositories {
		if !repository.IsNew {
			continue
		}
		impl, err := generateRepositoryImpl(repository.Repository)
		if err != nil {
			return "", err
		}
		src.WriteString("\n" + impl)
	}

	// Append new methods
	for _, repository := range repositories {
		for _, newMethod := range repository.NewMethods() {
			methodImpl, err := generateMethodImpl(repository.Repository, *newMethod)
			if err != nil {
				return "", err
			}
			src.WriteString("\n" + methodImpl)
		}
	}
	return formatImports(filepath, src.Bytes())
}

var repositoryStubFileTemplate = `
// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package {{ .Package }}
{{ range .MockDirectives -}}
//go:generate mockgen -source={{ .Src }} -destination={{ .Dst }}
{{ end -}}

{{ range .Imports }}
import {{ .Name }} "{{ .Path }}"
{{- end }}

var Repositories = fx.Options(
{{ range .Repositories -}}
  {{ .ImplPackage }}.{{ .QualifyString "Options" }},
{{ end -}}
)
`

func generateRepositoryStubFile(
	fsys fs.FS,
	packagePath string,
	repositories ...*RepositoryImpl,
) (string, error) {
	var templateData struct {
		Package        string
		Imports        []Import
		Repositories   []*RepositoryImpl
		MockDirectives []struct {
			Src string
			Dst string
		}
	}
	sort.Slice(repositories, func(i, j int) bool {
		a := repositories[i]
		b := repositories[j]
		if a.ImplPackage == b.ImplPackage {
			if a.Ident == "Repository" {
				return true
			}
			if b.Ident == "Repository" {
				return false
			}
			return a.Ident < b.Ident
		}
		return a.ImplPackage < b.ImplPackage
	})
	mocked := map[string]bool{}
	for _, repository := range repositories {
		src := path.Join(repository.PackagePath, repository.Filename)
		if _, done := mocked[src]; done {
			continue
		}
		mocked[src] = true
		dst := path.Join(repository.ImplPackagePath, "mocks", repository.Filename)
		templateData.MockDirectives = append(templateData.MockDirectives, struct {
			Src string
			Dst string
		}{
			Src: src,
			Dst: dst,
		})
	}

	templateData.Repositories = repositories
	pkgImport, pkgAlias, err := loadLocalPackage(fsys, nil, packagePath)
	if err != nil {
		return "", err
	}
	if pkgAlias != "" {
		templateData.Package = pkgAlias
	} else {
		templateData.Package = path.Base(pkgImport)
	}
	imports, err := collectImports(
		fsys,
		nil,
		false,
		true,
		repositories...,
	)
	if err != nil {
		return "", err
	}
	templateData.Imports = imports
	tmpl, err := template.
		New("repositoryStubFileTemplate").
		Parse(repositoryStubFileTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return formatImports(
		path.Join(packagePath, "repositories.go"),
		buf.Bytes(),
	)
}

func collectImports(
	fsys fs.FS,
	astFile *ast.File,
	importAPI, importImpl bool,
	repositories ...*RepositoryImpl,
) (allImports []Import, _ error) {
	usedImports := make(map[string]bool)
	if astFile != nil {
		for _, imp := range astFile.Imports {
			path, _ := strconv.Unquote(imp.Path.Value)
			usedImports[path] = true
		}
	}
	allImports = append(allImports, Import{Name: "", Path: "go.uber.org/fx"})
	if importAPI && importImpl {
		return nil, errors.New("cannot import both API and implementation")
	}
	for _, repository := range repositories {
		var rPkgPath string
		if importAPI {
			rPkgPath = repository.PackagePath
		} else {
			rPkgPath = repository.ImplPackagePath
		}
		rImport, rAlias, err := loadLocalPackage(
			fsys,
			astFile,
			rPkgPath,
		)
		if err != nil {
			return nil, err
		}
		// Check if there's a local package alias
		if astFile != nil && rAlias != "" {
			if importAPI {
				repository.Package = rAlias
			} else {
				repository.ImplPackage = rAlias
			}
		}
		allImports = append(allImports, Import{Name: rAlias, Path: rImport})
	}
	for _, repository := range repositories {
		allImports = append(allImports, repository.Imports...)
		for _, newMethod := range repository.NewMethods() {
			if newMethod.Params.HasCtx() {
				allImports = append(
					allImports,
					Import{Name: "", Path: "context"},
					Import{Name: "", Path: "go.opentelemetry.io/otel"},
					Import{Name: "", Path: "go.opentelemetry.io/otel/codes"},
				)
			}
			if newMethod.Returns.HasError() {
				allImports = append(allImports, Import{Name: "", Path: "github.com/rotisserie/eris"})
			}
		}
	}
	imports := []Import{}
	for _, imp := range allImports {
		if _, ok := usedImports[imp.Path]; !ok {
			imports = append(imports, imp)
		}
	}
	return imports, nil
}

func formatImports(filename string, src []byte) (string, error) {
	formattedSrc, err := imports.Process(filename, src, nil)
	if err != nil {
		return "", err
	}
	return string(formattedSrc), nil
}
