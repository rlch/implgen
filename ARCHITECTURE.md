# Architecture

This document provides a deep dive into the implgen codebase architecture, explaining how the different components work together to generate Go implementation files.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          CLI Layer                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │   main.go   │───▶│generate.go  │───▶│   urfave/cli/v3     │  │
│  │   Entry     │    │  Orchestrator│    │   Flag parsing     │  │
│  └─────────────┘    └─────────────┘    └─────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      File System Layer                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │   fs.go     │───▶│  crawlAPI   │───▶│  Path computation   │  │
│  │   Utils     │    │  Dir walker │    │  Module detection   │  │
│  └─────────────┘    └─────────────┘    └─────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Parser Layer                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │  parse.go   │───▶│ tree-sitter │───▶│   Query Engine      │  │
│  │  Interface  │    │   Parser    │    │   AST traversal     │  │
│  │  Extraction │    │             │    │   Pattern matching  │  │
│  └─────────────┘    └─────────────┘    └─────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Code Generation Layer                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │ codegen.go  │───▶│  Templates  │───▶│  File Generation    │  │
│  │ Generator   │    │   Engine    │    │  Import management  │  │
│  │             │    │             │    │  Formatting         │  │
│  └─────────────┘    └─────────────┘    └─────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. CLI Layer (`main.go`)

**Purpose**: Application entry point and command-line interface

**Key Responsibilities**:
- Parse command-line arguments using urfave/cli/v3
- Set up logging with structured output
- Validate user input
- Orchestrate the generation process

**Key Types**:
```go
var cmd = &cli.Command{
    Name:           "implgen",
    Description:    "Code generator for API implementations.",
    DefaultCommand: "generate",
    // ...
}
```

**Flow**:
1. Parse global flags (verbose, help)
2. Route to generate subcommand
3. Set up logger with appropriate level
4. Call generate function with parsed arguments

### 2. Generation Orchestrator (`generate.go`)

**Purpose**: Coordinates the overall generation process

**Key Responsibilities**:
- Crawl API directories for Go files
- Apply focus filters using glob patterns
- Coordinate parsing and code generation
- Handle file I/O operations

**Key Functions**:
```go
func generate(ctx context.Context, cmd *cli.Command) error
func groupByPackage(repositories []*RepositoryImpl) map[string][]*RepositoryImpl
func groupByImplFilename(repositories []*RepositoryImpl) map[string][]*RepositoryImpl
```

**Flow**:
1. Crawl API directory for `.go` files
2. Apply focus filters if specified
3. Parse repositories from each package
4. Compute implementation paths
5. Parse existing implementations
6. Generate new/updated implementation files
7. Generate dependency injection stub file

### 3. File System Layer (`fs.go`)

**Purpose**: Handle file system operations and path computations

**Key Responsibilities**:
- Walk directory trees to find Go files
- Compute implementation paths from API paths
- Detect Go module information
- Load local package imports

**Key Functions**:
```go
func crawlAPI(fsys fs.FS, apiDir string) (map[string][]string, error)
func computeImplPackagePath(apiRoot, implRoot, apiPackagePath string) (string, error)
func getModule(fsys fs.FS, root string) (string, error)
func loadLocalPackage(fsys fs.FS, astFile *ast.File, packagePath string) (string, string, error)
```

**Key Patterns**:
- Uses `fs.FS` interface for testability
- Caches module path for performance
- Handles both absolute and relative paths
- Integrates with Go's module system

### 4. Parser Layer (`parse.go`)

**Purpose**: Extract interface definitions from Go source code

**Key Technologies**:
- **tree-sitter**: Fast, robust parsing of Go syntax
- **go/parser**: Import extraction and validation
- **Custom queries**: S-expression queries for pattern matching

**Key Types**:
```go
type Repository struct {
    Package     string
    PackagePath string
    Filename    string
    Ident       string
    Generics    string
    Methods     []*Method
    Imports     []Import
}

type Method struct {
    Ident   string
    Params  Params
    Returns Params
}
```

**Tree-sitter Query**:
```scheme
(package_clause (package_identifier) @pkg) 

(type_spec
  name: (type_identifier) @class_name (#match? @class_name "Repository$")
  type_parameters: (type_parameter_list)? @generics 
  type: 
   (interface_type
     (method_elem
       name: (field_identifier) @method_name
       parameters: (parameter_list) @params
       result: [...]? @result)?))
```

**Parsing Flow**:
1. Parse imports using go/parser
2. Parse syntax tree using tree-sitter
3. Execute query to find Repository interfaces
4. Extract package name, interface names, generics
5. Parse method signatures and parameters
6. Build Repository structs with all metadata

### 5. Code Generation Layer (`codegen.go`)

**Purpose**: Generate implementation files from parsed interfaces

**Key Responsibilities**:
- Generate method implementations with proper signatures
- Add observability (OpenTelemetry) instrumentation
- Integrate error handling (eris wrapping)
- Generate dependency injection setup
- Manage imports and formatting

**Key Templates**:

#### Method Template
```go
const generateMethodTemplate = `
func (r *{{ .Repository.ImplName }}{{ .Repository.GenericsInstance }}) {{ .Method.Ident }}({{ .Method.Params.ParamsSrc }}){{ pad .Method.Returns.ReturnsSrc }}{
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
  {{- end }}
{{- end }}
  panic("TODO: implement {{ .Repository.QualifiedName }}.{{ .Method.Ident }}")
}
`
```

#### Repository Structure Template
```go
type {{ .Repository.QualifyString "Dependencies" }}{{ .Repository.Generics }} struct {
  fx.In  // or dig.In
  // Add dependencies here
}

func New{{ .Repository.Ident }}{{ .Repository.Generics }}(deps {{ .Repository.QualifyString "Dependencies" }}{{ .Repository.GenericsInstance }}) {{ .Repository.Package }}.{{ .Repository.Ident }}{{ .Repository.GenericsInstance }} {
  return &{{ .Repository.ImplName }}{{ .Repository.GenericsInstance }}{
    {{ .Repository.QualifyString "Dependencies" }}: deps,
  }
}

type {{ .Repository.ImplName }}{{ .Repository.Generics }} struct {
  {{ .Repository.QualifyString "Dependencies" }}{{ .Repository.GenericsInstance }}
}
```

**Generation Flow**:
1. Check for existing implementation files
2. Parse existing implementations to preserve them
3. Generate new repository structures for new interfaces
4. Generate method implementations for missing methods
5. Collect and organize imports
6. Format code with gofumpt and goimports
7. Write files to disk

## Data Flow

### Input Processing

```
API Directory
    │
    ├── package1/
    │   └── repository.go ──┐
    │                       │
    ├── package2/           │
    │   └── repository.go ──┼── crawlAPI() ──▶ map[string][]string
    │                       │                   (packagePath -> filenames)
    └── package3/           │
        └── repository.go ──┘
```

### Parsing Pipeline

```
Go Source Files
    │
    ├── tree-sitter parsing ──▶ Syntax Tree
    │                              │
    ├── go/parser ──▶ Import AST   │
    │                              │
    └────────────────┬─────────────┘
                     │
                     ▼
              Query Execution
                     │
                     ▼
              Repository structs
                     │
                     ▼
              RepositoryImpl structs
```

### Code Generation Pipeline

```
RepositoryImpl structs
    │
    ├── New repositories ──▶ Generate struct definitions
    │                         Generate constructor functions
    │                         Generate fx.Options/dig setup
    │
    └── Missing methods ──▶ Generate method implementations
                              Add observability code
                              Add error handling
                              Add TODO comments
                         │
                         ▼
                   Template execution
                         │
                         ▼
                   Import collection
                         │
                         ▼
                   Code formatting
                         │
                         ▼
                   File writing
```

## Key Algorithms

### 1. Type Qualification Algorithm

When generating method signatures, types need to be properly qualified:

```go
qualify = func(typ string) string {
    // Handle recursive types (maps, slices, pointers, functions)
    if strings.HasPrefix(typ, "map[") {
        // Recursively qualify key and value types
    } else if strings.HasPrefix(typ, "[") {
        // Handle slices and arrays
    } else if strings.HasPrefix(typ, "*") {
        // Handle pointers
    } else if strings.HasPrefix(typ, "func(") {
        // Handle function types
    }
    
    // Check if it's a built-in type or already qualified
    if strings.Contains(typ, ".") || isBuiltIn(typ) {
        return typ
    }
    
    // Check if it's a generic type parameter
    if slices.Contains(r.GenericsVariableList(), typ) {
        return typ
    }
    
    // Qualify with package name
    return r.Package + "." + typ
}
```

### 2. Import Collection Algorithm

Imports are collected from multiple sources and deduplicated:

```go
func collectImports(/* ... */) ([]Import, error) {
    usedImports := make(map[string]bool)
    
    // Add existing imports from AST
    for _, imp := range astFile.Imports {
        usedImports[imp.Path.Value] = true
    }
    
    // Add dependency injection import
    allImports = append(allImports, Import{Path: diPkgPath})
    
    // Add API package imports
    // Add implementation package imports
    
    // Add observability imports for methods with context
    // Add error handling imports for methods with errors
    
    // Deduplicate and return
    return deduplicatedImports, nil
}
```

### 3. File Merging Algorithm

When updating existing files, implementations are preserved:

```go
func generateRepositoryImplsForFile(/* ... */) (string, error) {
    // Read existing file if it exists
    if file != nil {
        // Parse existing content up to package declaration
        // Preserve existing imports and code
    }
    
    // Add new imports
    // Preserve existing file content
    
    // Append new repository structures
    for _, repository := range repositories {
        if repository.IsNew {
            // Generate new struct and constructor
        }
    }
    
    // Append new method implementations
    for _, repository := range repositories {
        for _, newMethod := range repository.NewMethods() {
            // Generate only missing methods
        }
    }
    
    return formatImports(filepath, src.Bytes())
}
```

## Dependency Management

### Tree-sitter Integration

The project uses a vendored copy of go-tree-sitter with the Go language parser:

```go
var (
    tsparser *sitter.Parser
    language *sitter.Language = sitter.NewLanguage(tsgo.Language())
)

func init() {
    tsparser = sitter.NewParser()
    err := tsparser.SetLanguage(language)
    if err != nil {
        panic(err)
    }
}
```

### Template Engine

Uses Go's `text/template` with custom functions:

```go
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
```

### Code Formatting

Two-stage formatting process:

1. **gofumpt**: Stricter formatting than gofmt
2. **goimports**: Import organization and addition

```go
func formatImports(filename string, src []byte) (string, error) {
    // First pass: gofumpt
    cmd := exec.Command("gofumpt")
    cmd.Stdin = bytes.NewReader(src)
    src, err = cmd.Output()
    
    // Second pass: goimports
    formattedSrc, err := imports.Process(filename, src, nil)
    return string(formattedSrc), nil
}
```

## Performance Considerations

### Caching

- **Module path caching**: Module path is cached after first detection
- **Tree-sitter parser reuse**: Single parser instance is reused
- **Query compilation**: Queries are compiled once and reused

### Memory Management

- **Tree cleanup**: tree-sitter trees are explicitly closed
- **Query cleanup**: Queries and cursors are properly disposed
- **Streaming processing**: Files are processed one at a time

### Parsing Efficiency

- **Targeted parsing**: Only interface definitions are extracted
- **Skip non-Go files**: Early filtering of file types
- **Import-only parsing**: Uses parser.ImportsOnly for existing files

## Error Handling

### Error Types

1. **File system errors**: Missing directories, permission issues
2. **Parsing errors**: Invalid Go syntax, tree-sitter failures
3. **Generation errors**: Template execution, formatting failures
4. **Module errors**: Invalid go.mod, missing module path

### Error Wrapping

Consistent error wrapping with context:

```go
if err != nil {
    return fmt.Errorf("failed to parse repositories in %s: %w", apiPackagePath, err)
}
```

### Graceful Degradation

- **Skip invalid files**: Continue processing other files
- **Preserve existing code**: Never overwrite working implementations
- **Detailed logging**: Verbose mode provides debugging information

## Testing Strategy

### Unit Tests

- **Parser tests**: Verify interface extraction from various Go syntax
- **Generator tests**: Validate template execution and output
- **File system tests**: Mock file system operations

### Integration Tests

- **End-to-end tests**: Full generation pipeline with example projects
- **Regression tests**: Ensure backward compatibility
- **Cross-platform tests**: Verify behavior on different operating systems

### Test Data

The `example/` directories serve as comprehensive test cases, covering:
- Simple interfaces
- Generic interfaces
- Nested packages
- Multiple repositories per package
- Both fx and dig patterns

## Extensibility

### Adding New DI Frameworks

To support a new dependency injection framework:

1. Add new flag to CLI
2. Modify template selection in `generateRepositoryImpl`
3. Update import collection logic
4. Add new stub template in `generateRepositoryStubFile`

### Supporting New Interface Patterns

To support interfaces with different naming patterns:

1. Modify tree-sitter query in `parseRepositories`
2. Update interface detection logic
3. Ensure backward compatibility

### Adding New Code Generation Features

To add new generated code features:

1. Extend Repository/Method structs with new metadata
2. Update parsing logic to extract new information
3. Modify templates to include new features
4. Update import collection as needed

This architecture provides a solid foundation for the current feature set while remaining extensible for future enhancements.