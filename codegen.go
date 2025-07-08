// Package main contains code generation logic for creating Go implementation files.
//
// This file implements the template-based code generation system that creates
// implementation files with proper method signatures, dependency injection setup,
// observability instrumentation, and error handling. It manages imports, formatting,
// and merging with existing implementation files.
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
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

// ImplTestPackage returns the test package name for the implementation package.
// This follows Go's convention of using "_test" suffix for external test packages.
func (c ContractImpl) ImplTestPackage() string {
	return c.ImplPackage + "_test"
}

// NewMethods returns the list of methods that need to be implemented.
//
// It compares the interface methods against the existing implementation methods
// and returns only those that are missing. The returned methods have their
// parameter and return types properly qualified with package names.
func (c ContractImpl) NewMethods() []*Method {
	methods := []*Method{}
	for _, method := range c.Methods {
		existing := slices.Contains(c.ImplMethods, method.Ident)
		if existing {
			continue
		}
		args := make(Params, len(method.Params))
		returns := make(Params, len(method.Returns))

		var qualify func(string) string
		qualify = func(typ string) string {
			n := len(typ)
			typ = strings.TrimSpace(typ)
			if typ == "" {
				return ""
			}
			// Handle recursive types
			if strings.HasPrefix(typ, "map[") {
				start, end := getEnclosingBrackets(typ, '[', ']')
				return "map[" + qualify(typ[start+1:end]) + "]" + qualify(typ[end+1:])
			} else if strings.HasPrefix(typ, "[") {
				splitIdx := strings.Index(typ, "]")
				return typ[0:splitIdx+1] + qualify(typ[splitIdx+1:])
			} else if strings.HasPrefix(typ, "*") {
				return "*" + qualify(typ[1:])
			} else if strings.HasPrefix(typ, "...") {
				return "..." + qualify(typ[3:])
			} else if strings.HasPrefix(typ, "func(") {
				start, end := getEnclosingBrackets(typ, '(', ')')
				args := parseParams(typ[start+1 : end])
				returns := parseParams(typ[end+1:])
				for _, p := range append(args, returns...) {
					p.Type = qualify(p.Type)
				}
				return "func(" + args.ParamsSrc() + ") " + returns.ReturnsSrc()
			} else if genericStart := strings.Index(typ, "["); genericStart != -1 {
				// We know the last character is a ] as it's a generic and have handled
				// other composite types above.
				genericVars := strings.Split(typ[genericStart+1:n-1], ",")
				for i, g := range genericVars {
					genericVars[i] = qualify(strings.TrimSpace(g))
				}
				return qualify(typ[:genericStart]) + "[" + strings.Join(genericVars, ", ") + "]"
			}
			if len(typ) == 0 {
				return ""
			}
			isLower := 'a' <= typ[0] && typ[0] <= 'z'
			// . implies package qualification, lowercase implies built-in
			if strings.Contains(typ, ".") || isLower {
				return typ
			}
			// handle case where we return a generic defined by repository
			if slices.Contains(c.GenericsVariableList(), typ) {
				return typ
			}
			return c.Package + "." + typ
		}
		for i, arg := range method.Params {
			arg := arg
			args[i] = &Param{
				Ident: arg.Ident,
				Type:  qualify(arg.Type),
			}
		}
		for i, arg := range method.Returns {
			arg := arg
			returns[i] = &Param{
				Ident: arg.Ident,
				Type:  qualify(arg.Type),
			}
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

func (p Params) QualifyNames() {
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
	p.QualifyNames()
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
	p.QualifyNames()
	src := p.ParamsSrc()
	if p.Named() || len(p) > 1 {
		return "(" + src + ")"
	}
	return src
}

func (c Contract) QualifyString(s string) string {
	name := c.Name()
	if name == "Repository" {
		return s
	}
	return name + s
}

func (c Contract) Name() string {
	// Extract the base name by removing the suffix
	// This works for any suffix (Repository, Service, Handler, etc.)
	for _, suffix := range []string{"Repository", "Service", "Handler", "Manager", "Controller"} {
		if strings.HasSuffix(c.Ident, suffix) {
			return c.Ident[:len(c.Ident)-len(suffix)]
		}
	}
	return c.Ident
}

func (c Contract) ImplName() string {
	if c.Ident == "" {
		return ""
	}
	name := c.Ident
	return strings.ToLower(string(name[0])) + name[1:] + "Impl"
}

func (c Contract) QualifiedName() string {
	if c.Package == "" {
		return c.Ident
	}
	return c.Package + "." + c.Ident
}

const generateMethodTemplate = `
  func (r *{{ .Contract.ImplName }}{{ .Contract.GenericsInstance }}) {{ .Method.Ident }}({{ .Method.Params.ParamsSrc }}){{ pad .Method.Returns.ReturnsSrc }}{
  {{- if .Method.Params.HasCtx }}
    ctx, span := otel.GetTracerProvider().Tracer("{{ .Contract.Package }}").Start(ctx, "{{ .Contract.Name }}.{{ .Method.Ident }}")
    {{- if .Method.Returns.HasError }}
    defer func() {
      if err != nil {
        err = fault.Wrap(err, fmsg.With("{{ .Contract.QualifiedName }}.{{ .Method.Ident }}"))
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
        err = fault.Wrap(err, fmsg.With("{{ .Contract.QualifiedName }}.{{ .Method.Ident }}"))
      }
    }()
    {{- end }}
  {{- end }}
    panic("TODO: implement {{ .Contract.QualifiedName }}.{{ .Method.Ident }}")
  }
`

func generateMethodImpl(contract Contract, method Method) (string, error) {
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
		Contract
		Method
	}{contract, method}); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// generateContractImpl generates the method and struct declarations for a single contract.
func generateContractImpl(contract Contract) (string, error) {
	diPkg := "fx"
	options := `
var {{ .Contract.QualifyString "Options" }} = fx.Options(
	fx.Provide(
		New{{ .Contract.Ident }},
	),
)
`
	if fUseDig {
		diPkg = "dig"
		options = ""
	} else if contract.Generics != "" {
		options = ""
	}
	generateContractImplTemplate := `
type {{ .Contract.QualifyString "Dependencies" }}{{ .Contract.Generics }} struct {
  ` + diPkg + `.In
	// Add dependencies here
}
` + options + `
func New{{ .Contract.Ident }}{{ .Contract.Generics }}(deps {{ .Contract.QualifyString "Dependencies" }}{{ .Contract.GenericsInstance }}) {{ .Contract.Package }}.{{ .Contract.Ident }}{{ .Contract.GenericsInstance }} {
	return &{{ .Contract.ImplName }}{{ .Contract.GenericsInstance }}{
    {{ .Contract.QualifyString "Dependencies" }}: deps,
	}
}

type {{ .Contract.ImplName }}{{ .Contract.Generics }} struct {
  {{ .Contract.QualifyString "Dependencies" }}{{ .Contract.GenericsInstance }}
}
`
	tmpl, err := template.
		New("generateContractImplTemplate").
		Parse(generateContractImplTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Contract
	}{contract}); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// generateContractImplsForFile generates the contract implementations for a single file.
//
// All ContractImpl's are assumed to be for the same file as contracts[0].
func generateContractImplsForFile(
	fsys fs.FS,
	filepath string,
	contracts []*ContractImpl,
) (_ string, err error) {
	if len(contracts) == 0 {
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
		defer func() { _ = file.Close() }()
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
// This file will be automatically regenerated based on the API. Any contract implementations
// will be copied through when generating and new methods will be added to the end.
package %s
`, contracts[0].ImplPackage)
		src.WriteString(strings.TrimPrefix(packageDecl, "\n"))
	}

	// Add imports to src
	requiredImports, err := collectImports(
		fsys,
		astFile,
		true,
		false,
		nil,
		contracts...,
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

	// Append new contract declarations
	for _, contract := range contracts {
		if !contract.IsNew {
			continue
		}
		impl, err := generateContractImpl(contract.Contract)
		if err != nil {
			return "", err
		}
		src.WriteString("\n" + impl)
	}

	// Append new methods
	for _, contract := range contracts {
		for _, newMethod := range contract.NewMethods() {
			methodImpl, err := generateMethodImpl(contract.Contract, *newMethod)
			if err != nil {
				return "", err
			}
			src.WriteString("\n" + methodImpl)
		}
	}
	return formatImports(filepath, src.Bytes())
}

func generateContractStubFile(
	fsys fs.FS,
	packagePath string,
	contracts ...*ContractImpl,
) (string, error) {
	type MockDirective struct {
		Src          string
		Dst          string
		ImplPackage  string
		Contracts []string
	}
	var templateData struct {
		Package        string
		Imports        []Import
		Contracts   []*ContractImpl
		MockDirectives []MockDirective
	}
	sort.Slice(contracts, func(i, j int) bool {
		a := contracts[i]
		b := contracts[j]
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
	for _, contractGroup := range groupByPackage(contracts) {
		contract := contractGroup[0]
		src := contract.PackagePath
		if _, done := mocked[src]; done {
			continue
		}
		var err error
		src, err = filepath.Rel(fImpl, src)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
		mocked[src] = true
		dst := path.Join(contract.ImplPackagePath, "mocks.go")
		dst, err = filepath.Rel(fImpl, dst)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
		contractIdents := make([]string, len(contractGroup))
		for i, contract := range contractGroup {
			contractIdents[i] = contract.Ident
		}
		slices.Sort(contractIdents)
		templateData.MockDirectives = append(templateData.MockDirectives, MockDirective{
			Src:          src,
			Dst:          dst,
			ImplPackage:  contract.ImplPackage,
			Contracts: contractIdents,
		})
	}
	sort.Slice(templateData.MockDirectives, func(i, j int) bool {
		return templateData.MockDirectives[i].ImplPackage < templateData.MockDirectives[j].ImplPackage
	})

	templateData.Contracts = contracts
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
		nil,
		contracts...,
	)
	if err != nil {
		return "", err
	}
	templateData.Imports = imports

	var repositoryStubFileTemplate string
	if fUseDig {
		repositoryStubFileTemplate = `
// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package {{ .Package }}
{{ range .MockDirectives -}}
//go:generate moq -out={{ .Dst }} -pkg={{ .ImplPackage }} -rm -skip-ensure {{ .Src }} {{ range .Contracts }}{{.}} {{ end }}
{{ end -}}

import (
{{- range .Imports }}
	{{ if .Name }}{{ .Name }} {{ end }}"{{ .Path }}"
{{- end }}

	_ "github.com/Southclaws/fault"
	_ "github.com/Southclaws/fault/fmsg"
)

var RepositoryFactories = []any{
{{ range .Contracts -}}
  {{ if not .Generics -}} 
  {{ .ImplPackage }}.New{{ .Ident }},
  {{- end }}
{{ end -}}
}
`
	} else {
		repositoryStubFileTemplate = `
// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package {{ .Package }}
{{ range .MockDirectives -}}
//go:generate moq -out={{ .Dst }} -pkg={{ .ImplPackage }} -rm -skip-ensure {{ .Src }} {{ range .Contracts }}{{.}} {{ end }}
{{ end -}}

import (
{{- range .Imports }}
	{{ if .Name }}{{ .Name }} {{ end }}"{{ .Path }}"
{{- end }}

	_ "github.com/Southclaws/fault"
	_ "github.com/Southclaws/fault/fmsg"
)

var Repositories = fx.Options(
{{ range .Contracts -}}
  {{ if not .Generics -}} 
  {{ .ImplPackage }}.{{ .QualifyString "Options" }},
  {{- end }}
{{ end -}}
)
`
	}
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
	extraImports []Import,
	contracts ...*ContractImpl,
) (allImports []Import, _ error) {
	usedImports := make(map[string]bool)
	if astFile != nil {
		for _, imp := range astFile.Imports {
			path, _ := strconv.Unquote(imp.Path.Value)
			usedImports[path] = true
		}
	}
	diPkgPath := "go.uber.org/fx"
	if fUseDig {
		diPkgPath = "go.uber.org/dig"
	}
	allImports = append(allImports, Import{Name: "", Path: diPkgPath})
	allImports = append(allImports, extraImports...)
	importContracts := func(importAPI, importImpl bool) error {
		for _, contract := range contracts {
			var rPkgPath string
			if importAPI {
				rPkgPath = contract.PackagePath
			} else if importImpl {
				rPkgPath = contract.ImplPackagePath
			}
			rImport, rAlias, err := loadLocalPackage(
				fsys,
				astFile,
				rPkgPath,
			)
			if err != nil {
				return err
			}
			// Check if there's a local package alias
			if astFile != nil && rAlias != "" {
				if importAPI {
					contract.Package = rAlias
				} else if importImpl {
					contract.ImplPackage = rAlias
				}
			}
			// For implementation packages, ensure we have an alias
			if importImpl && rAlias == "" {
				rAlias = contract.ImplPackage
			}
			allImports = append(
				allImports,
				Import{
					Name: rAlias,
					Path: rImport,
				},
			)
		}
		return nil
	}
	if err := importContracts(importAPI, false); err != nil {
		return nil, err
	}
	if err := importContracts(false, importImpl); err != nil {
		return nil, err
	}

	for _, contract := range contracts {
		allImports = append(allImports, contract.Imports...)
		for _, newMethod := range contract.NewMethods() {
			if newMethod.Params.HasCtx() {
				allImports = append(
					allImports,
					Import{Name: "", Path: "context"},
					Import{Name: "", Path: "go.opentelemetry.io/otel"},
					Import{Name: "", Path: "go.opentelemetry.io/otel/codes"},
				)
			}
			if newMethod.Returns.HasError() {
				allImports = append(allImports, 
					Import{Name: "", Path: "github.com/Southclaws/fault"},
					Import{Name: "", Path: "github.com/Southclaws/fault/fmsg"},
				)
			}
		}
	}
	imports := []Import{}
	for _, imp := range allImports {
		if _, ok := usedImports[imp.Path]; !ok {
			usedImports[imp.Path] = true
			imports = append(imports, imp)
		}
	}
	return imports, nil
}

func formatImports(filename string, src []byte) (_ string, err error) {
	cmd := exec.Command("gofumpt")
	cmd.Stdin = bytes.NewReader(src)
	src, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run gofumpt: %w", err)
	}
	formattedSrc, err := imports.Process(filename, src, nil)
	if err != nil {
		return "", fmt.Errorf("failed to process imports: %w", err)
	}
	return string(formattedSrc), nil
}
