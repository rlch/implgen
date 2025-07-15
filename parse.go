// Package main contains parsing logic for extracting Go interface definitions.
//
// This file implements the core parsing functionality that uses tree-sitter to analyze
// Go source code and extract Repository interface definitions along with their methods,
// parameters, and type information. It also handles parsing of existing implementation
// files to determine what methods already exist.
package main

import (
	"errors"
	"fmt"
	"go/parser"
	"io"
	"io/fs"
	"log/slog"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/danielgtaylor/casing"
	tsgo "github.com/rlch/implgen/parser/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var (
	// ErrNoPackage is returned when no package declaration is found in a Go file.
	ErrNoPackage = errors.New("no package name found")

	// tsparser is the global tree-sitter parser instance used for parsing Go source code.
	tsparser *sitter.Parser

	// language represents the Go language grammar for tree-sitter parsing.
	language *sitter.Language = sitter.NewLanguage(tsgo.Language())
)

func init() {
	tsparser = sitter.NewParser()
	err := tsparser.SetLanguage(language)
	if err != nil {
		panic(err)
	}
}

type (
	// Contract represents a parsed Go interface with a configurable suffix.
	// It contains all the metadata needed to generate an implementation.
	Contract struct {
		Package       string    // Package name (e.g., "user")
		PackagePath   string    // Full path to the package directory
		Filename      string    // Name of the file containing the interface
		Ident         string    // Interface identifier (e.g., "UserRepository", "UserService")
		Generics      string    // Generic type parameters (e.g., "[T any]")
		Methods       []*Method // All methods defined in the interface
		Imports       []Import  // Import statements from the source file
		Embeds        []string  // Embedded interfaces in this contract
		IgnoredEmbeds []string  // Embedded interfaces marked with //implgen:ignore
		Ignored       bool      // Whether entire contract is marked with //implgen:ignore
	}

	// ContractImpl extends Contract with implementation-specific metadata.
	// It tracks what implementations already exist and where they should be generated.
	ContractImpl struct {
		Contract
		IsNew           bool     // Whether this is a completely new implementation
		ImplPackage     string   // Implementation package name (e.g., "userimpl")
		ImplPackagePath string   // Path to implementation package directory
		ImplFilename    string   // Name of the implementation file
		ImplMethods     []string // Names of methods that already have implementations
	}

	// Import represents a Go import statement.
	Import struct {
		Name string // Import alias (empty for no alias)
		Path string // Import path
	}

	// Method represents a single method in an interface.
	Method struct {
		Ident   string // Method name
		Params  Params // Method parameters
		Returns Params // Method return values
		Ignored bool   // Whether method is marked with //implgen:ignore
	}

	// Params is a slice of parameters or return values.
	Params []*Param

	// Param represents a single parameter or return value.
	Param struct {
		Ident string // Parameter name (may be empty for unnamed parameters)
		Type  string // Parameter type
	}

)

// GenericsVariableList extracts the generic type variable names from the generics string.
// For example, "[T any, U comparable]" returns ["T", "U"].
func (c Contract) GenericsVariableList() []string {
	out := []string{}
	generics := strings.Trim(c.Generics, "[]")
	if generics == "" {
		return nil
	}
	for _, s := range strings.Split(generics, ",") {
		s = strings.TrimSpace(s)
		out = append(out, strings.Split(s, " ")[0])
	}
	return out
}

// GenericsInstance returns the generic type instantiation string for use in implementations.
// For example, "[T any, U comparable]" becomes "[T, U]" for instantiating the generic type.
func (c Contract) GenericsInstance() (out string) {
	generics := strings.Trim(c.Generics, "[]")
	if generics == "" {
		return ""
	}
	for i, s := range strings.Split(generics, ",") {
		if i != 0 {
			out += ", "
		}
		s = strings.TrimSpace(s)
		out += strings.Split(s, " ")[0]
	}
	return "[" + out + "]"
}

// parseContractsForPackage extracts Contract interfaces from all Go files in a package.
//
// It processes each file using tree-sitter to parse the syntax tree and extract
// interface definitions that end with the specified suffix. The function aggregates
// contracts from all files in the package and sets their package metadata.
func parseContractsForPackage(
	fsys fs.FS,
	packagePath string,
	packageFiles []string,
	suffix string,
) (contracts []*Contract, err error) {
	contracts = []*Contract{}
	for _, filename := range packageFiles {
		fullPath := path.Join(packagePath, filename)
		file, err := fsys.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fullPath, err)
		}
		defer func() { _ = file.Close() }()
		var src []byte
		src, err = io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
		}
		tree := tsparser.Parse(src, nil)
		defer tree.Close()
		packageContracts, err := parseContracts(src, tree, suffix)
		if err != nil {
			return nil, fmt.Errorf("failed to extract contracts from file %s: %w", fullPath, err)
		}
		for _, contract := range packageContracts {
			contract.Filename = filename
			contract.PackagePath = packagePath
		}
		contracts = append(contracts, packageContracts...)
	}

	filtered := contracts[:0]
	for _, c := range contracts {
		if !c.Ignored {
			filtered = append(filtered, c)
		}
	}
	contracts = filtered
	return contracts, nil
}

// resolveEmbeddedInterfaces recursively resolves methods from embedded interfaces
func resolveEmbeddedInterfaces(contracts []*Contract) {
	// Create a map for quick lookup of contracts by name
	contractMap := make(map[string]*Contract)
	for _, contract := range contracts {
		contractMap[contract.Ident] = contract
	}
	
	// Process each contract
	for _, contract := range contracts {
		resolveEmbeddedInterfacesForContract(contract, contractMap, make(map[string]bool))
	}
}

// resolveEmbeddedInterfacesForContract resolves embedded interfaces for a single contract
func resolveEmbeddedInterfacesForContract(contract *Contract, contractMap map[string]*Contract, visited map[string]bool) {
	// Prevent infinite recursion
	if visited[contract.Ident] {
		return
	}
	visited[contract.Ident] = true
	
	// Process each embedded interface
	for _, embedName := range contract.Embeds {
		if embeddedContract, exists := contractMap[embedName]; exists {
			if embeddedContract.Ignored {
				continue
			}
			// First resolve the embedded contract's own embedded interfaces
			resolveEmbeddedInterfacesForContract(embeddedContract, contractMap, visited)
			
			// Add methods from the embedded interface
			for _, method := range embeddedContract.Methods {
				// Check if method already exists to avoid duplicates
				exists := false
				for _, existingMethod := range contract.Methods {
					if existingMethod.Ident == method.Ident {
						exists = true
						break
					}
				}
				if !exists {
					// Create a copy of the method
					methodCopy := &Method{
						Ident:   method.Ident,
						Ignored: method.Ignored,
					}
					
					// Copy params (handle nil case)
					if method.Params != nil {
						methodCopy.Params = make([]*Param, len(method.Params))
						copy(methodCopy.Params, method.Params)
					}
					
					// Copy returns (handle nil case)
					if method.Returns != nil {
						methodCopy.Returns = make([]*Param, len(method.Returns))
						copy(methodCopy.Returns, method.Returns)
					}
					
					contract.Methods = append(contract.Methods, methodCopy)
				}
			}
		}
	}
}

func parseContracts(src []byte, tree *sitter.Tree, suffix string) (contracts []*Contract, err error) {
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
		for _, contract := range contracts {
			contract.Imports = imports
		}
	}()

	const (
		PKG_CAPTURE = "pkg"
		CLASS_NAME_CAPTURE = "class_name"
		GENERICS_CAPTURE = "generics"
		METHOD_NAME_CAPTURE = "method_name"
		PARAMS_CAPTURE = "params"
		RESULT_CAPTURE = "result"
		TYPE_COMMENT_CAPTURE = "type_comment"
		METHOD_COMMENT_CAPTURE = "method_comment"
		EMBED_NAME_CAPTURE = "embed_name"
	)
	// Build dynamic query with configurable suffix
	queryString := fmt.Sprintf(`
(package_clause (package_identifier) @pkg) 

(comment) @type_comment

(comment) @method_comment

(type_spec
  name: (type_identifier) @class_name (#match? @class_name "%s$")
  type_parameters: (type_parameter_list)? @generics 
  type: 
   (interface_type
     (method_elem
       name: (field_identifier) @method_name
       parameters: (parameter_list) @params
       result: [
        (parameter_list)
        (type_identifier)
        (qualified_type)
       ]? @result)?
     (type_elem
       (type_identifier) @embed_name)?))
    `, suffix)
	query, queryErr := sitter.NewQuery(language, queryString)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to create query: %s", queryErr)
	}
	defer query.Close()
	
	// Get capture names for string-based matching
	captureNames := query.CaptureNames()
	
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()
	qc := cursor.Captures(query, tree.RootNode(), src)

	// Get package name
	m, _ := qc.Next()
	if m == nil {
		return nil, nil
	}
	if len(m.Captures) != 1 {
		return nil, ErrNoPackage
	}

	pkg := m.Captures[0].Node.Utf8Text(src)
	defer func() {
		for _, contract := range contracts {
			contract.Package = pkg
		}
	}()
	m, _ = qc.Next()
	if m == nil {
		return nil, nil
	}
	var curIdx, methodIdx int
	var pendingTypeComment, pendingMethodComment string
	contracts = append(contracts, &Contract{})
	for {
		for _, c := range m.Captures {
			contract := contracts[curIdx]
			nodeSrc := c.Node.Utf8Text(src)
			captureName := captureNames[c.Index]
			switch captureName {
			case CLASS_NAME_CAPTURE:
				name := nodeSrc
				// NOTE: Apparently there's a problem with #match? directive, hacky
				// workaround
				if !strings.HasSuffix(name, suffix) {
					continue
				}
				if contract.Ident == name {
					continue
				} else if contract.Ident != "" {
					// We have a new contract, so append the previous one to the list.
					methodIdx = 0
					curIdx++
					contracts = append(contracts, &Contract{Ident: name})
					contract = contracts[curIdx]
				} else if contract.Ident == "" {
					contract.Ident = name
				}
				
				// Apply pending type comment if it contains ignore directive
				if pendingTypeComment != "" && strings.Contains(pendingTypeComment, "implgen:ignore") {
					contract.Ignored = true
				}
				pendingTypeComment = "" // Clear after processing
				
				// Also clear method comment to prevent it from being applied to methods in an ignored contract
				if contract.Ignored {
					pendingMethodComment = ""
				}
			case GENERICS_CAPTURE:
				contract.Generics = nodeSrc
			case METHOD_NAME_CAPTURE:
				var curMethod *Method
				methodName := nodeSrc
				if len(contract.Methods) == 0 {
					curMethod = &Method{Ident: methodName}
					contract.Methods = append(contract.Methods, curMethod)
				} else {
					curMethod = contract.Methods[methodIdx]
				}
				if curMethod.Ident != methodName {
					found := false
					for i, m := range contract.Methods {
						if m.Ident == methodName {
							methodIdx = i
							found = true
							break
						}
					}
					if !found {
						methodIdx = len(contract.Methods)
						curMethod = &Method{Ident: methodName}
						contract.Methods = append(contract.Methods, curMethod)
					}
				}
				
				// Apply pending method comment if it contains ignore directive
				if pendingMethodComment != "" && strings.Contains(pendingMethodComment, "implgen:ignore") {
					contract.Methods[methodIdx].Ignored = true
				}
				pendingMethodComment = "" // Clear after processing
			case PARAMS_CAPTURE:
				contract.Methods[methodIdx].Params = parseParams(nodeSrc)
			case RESULT_CAPTURE:
				contract.Methods[methodIdx].Returns = parseParams(nodeSrc)
			case TYPE_COMMENT_CAPTURE:
				// Store comment for potential application to next contract
				// Only store if it's an ignore comment to avoid capturing unrelated comments
				if strings.Contains(nodeSrc, "implgen:ignore") {
					pendingTypeComment = nodeSrc
				}
			case METHOD_COMMENT_CAPTURE:
				// Store comment for potential application to next method
				// Only store if it's an ignore comment to avoid capturing unrelated comments
				if strings.Contains(nodeSrc, "implgen:ignore") {
					pendingMethodComment = nodeSrc
				}
			case EMBED_NAME_CAPTURE:
				embedName := nodeSrc
				// Add embedded interface to current contract if not already present
				found := false
				for _, existing := range contract.Embeds {
					if existing == embedName {
						found = true
						break
					}
				}
				if !found {
					contract.Embeds = append(contract.Embeds, embedName)
				}
			default:
				slog.Error(
					"unhandled",
					slog.Int("index", int(c.Index)),
					slog.String("src", c.Node.Utf8Text(src)),
				)
			}
		}
		m, _ = qc.Next()
		if m == nil {
			break
		}
	}
	if len(contracts) == 1 && contracts[0].Ident == "" {
		return nil, nil
	}
	
	// Resolve embedded interfaces
	resolveEmbeddedInterfaces(contracts)
	
	return
}

func getEnclosingBrackets(s string, left, right rune) (start, end int) {
	bracketCount := 0
	start = strings.Index(s, string(left))
	for i, c := range s[start+1:] {
		switch c {
		case left:
			bracketCount++
		case right:
			bracketCount--
		}
		if bracketCount == -1 {
			return start, start + 1 + i
		}
	}
	return -1, -1
}

func parseParams(src string) Params {
	src = strings.TrimSpace(src)
	if src == "" {
		return nil
	}
	if src[0] == '(' {
		src = src[1 : len(src)-1]
	}
	// We need to handle arguments accepting a comma so can't just split on a
	// comma.
	args := []string{}
	var (
		lastComma = -1
		i         int
	)
	for {
		if i >= len(src) {
			// Single argument
			if lastComma == -1 {
				args = append(args, src)
			} else {
				lastArg := strings.TrimSpace(src[lastComma+1:])
				if lastArg != "" {
					args = append(args, lastArg)
				}
			}
			break
		}
		c := src[i]
		if c == '(' {
			_, end := getEnclosingBrackets(src[i:], '(', ')')
			i += end + 1
			continue
		}
		if c == '[' {
			_, end := getEnclosingBrackets(src[i:], '[', ']')
			i += end + 1
			continue
		}
		if c == ',' {
			args = append(args, src[lastComma+1:i])
			lastComma = i
		}
		i++
	}
	if len(args) == 0 || (len(args) == 1 && args[0] == "") {
		return nil
	}
	named := false
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		parts := strings.Split(arg, " ")
		if strings.HasPrefix(parts[0], "func(") {
			break
		}
		if len(parts) > 1 {
			named = true
			break
		}
	}
	if named {
		return parseNamedParams(args)
	}
	params := make([]*Param, len(args))
	for i, part := range args {
		params[i] = &Param{Type: strings.TrimSpace(part)}
	}
	return params
}

func parseNamedParams(parts []string) Params {
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

func parseContractImpls(
	fsys fs.FS,
	implPackagePath string,
	contracts []*Contract,
) ([]*ContractImpl, error) {
	if len(contracts) == 0 {
		return nil, nil
	}

	defaultImplFilename := func(contract *ContractImpl) string {
		return casing.Snake(contract.Name()) + "_impl.go"
	}
	implPackageName := contracts[0].Package + "impl"
	impls := make([]*ContractImpl, len(contracts))
	for i, contract := range contracts {
		impls[i] = &ContractImpl{
			Contract: *contract,
		}
	}
	entries, err := fs.ReadDir(fsys, implPackagePath)
	if err != nil {
		err := err.(*fs.PathError)
		if errors.Is(err.Err, fs.ErrNotExist) {
			for _, contract := range impls {
				contract.ImplPackage = implPackageName
				contract.ImplPackagePath = implPackagePath
				contract.ImplFilename = defaultImplFilename(contract)
				contract.IsNew = true
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
		defer func() { _ = file.Close() }()
		src, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}
		pkg, implDecls, methods, err := parseRepositoryImplFile(src)
		if err == ErrNoPackage {
			continue
		} else if err != nil {
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

	for _, contract := range impls {
		implName := contract.ImplName()
		contract.ImplPackage = implPackageName
		contract.ImplPackagePath = implPackagePath
		if filename, ok := implDeclsToFileMap[implName]; ok {
			contract.ImplFilename = filename
			contract.ImplMethods = repositoryToMethodMap[implName]
		} else {
			contract.ImplFilename = defaultImplFilename(contract)
			contract.IsNew = true
		}
	}
	return impls, nil
}

func parseRepositoryImplFile(src []byte) (
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
	// NOTE: We use (.*) after an Impl as a workaround for (\[.*\])?
	// should be fiiiiiiiiiiiiiiiiiiiiinee
	query, queryErr := sitter.NewQuery(language, `
  (
    (package_clause (package_identifier) @pkg)
    (type_declaration 
      (type_spec
          name: (type_identifier) @impl_name (#match? @impl_name "Impl$")))?
    (method_declaration
        receiver: (parameter_list
          (parameter_declaration
            type: (_) @impl_rec (#match? @impl_rec "Impl(.*)?$")))
        name: (field_identifier) @impl_field)?
  )
`)
	if queryErr != nil {
		return "", nil, nil, queryErr
	}
	defer query.Close()

	tree := tsparser.Parse(src, nil)
	defer tree.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()
	qc := cursor.Captures(query, tree.RootNode(), src)

	repImplMap := map[string]bool{}
	var curRec string
	for {
		m, _ := qc.Next()
		if m == nil {
			break
		}
		for _, c := range m.Captures {
			nodeSrc := c.Node.Utf8Text(src)
			switch c.Index {
			case PKG_CAPTURE:
				packageName = nodeSrc
			case IMPL_NAME_CAPTURE:
				implName := nodeSrc
				repImplMap[implName] = true
			case IMPL_REC_CAPTURE:
				rec := nodeSrc
				if rec == "" {
					continue
				}
				if rec[0] == '*' {
					rec = rec[1:]
				}
				if idx := strings.Index(rec, "["); idx != -1 {
					rec = rec[:idx]
				}
				curRec = rec
			case IMPL_FIELD_CAPTURE:
				if curRec == "" {
					panic("receiver not found")
				}
				method := nodeSrc
				found := slices.Contains(methods[curRec], method)
				if !found {
					methods[curRec] = append(methods[curRec], method)
				}
			default:
				slog.Error(
					"unhandled",
					slog.Int("index", int(c.Index)),
					slog.String("src", nodeSrc),
				)
			}
		}
	}
	if packageName == "" {
		return "", nil, nil, ErrNoPackage
	}
	for impl := range repImplMap {
		repImpls = append(repImpls, impl)
	}
	return
}
