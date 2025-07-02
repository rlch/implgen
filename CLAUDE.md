# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`implgen` is a Go code generator that automatically creates implementation boilerplate for Go interfaces. It parses interface definitions (typically ending with "Repository") from an `api` directory and generates corresponding implementation files in an `internal` directory.

## Key Commands

### Build
```bash
go build -o implgen .
```

### Test
```bash
go test ./...
```
Note: Tests will fail if tree-sitter dependencies are not properly built. Focus on the core implgen functionality tests.

### Lint
```bash
golangci-lint run --config .golangci.yaml ./...
```

### Generate Implementation Files
```bash
# Generate all implementations
./implgen generate

# With custom directories
./implgen generate --root . --api api --impl internal

# Focus on specific packages (using glob patterns)
./implgen generate --focus "waltuh/**"

# Use dig instead of fx for dependency injection
./implgen generate --dig
```

## Architecture

### Core Components

1. **Parser (`parse.go`)**: Uses go-tree-sitter to parse Go source files and extract interface definitions. The parser identifies interfaces ending with "Repository" and extracts their methods, parameters, and return types.

2. **Code Generator (`codegen.go`)**: Generates implementation files with:
   - Proper package structure
   - Import statements
   - Struct definitions with embedded dependencies
   - Method stubs with panic("TODO: implement...") 
   - OpenTelemetry tracing for methods with context
   - Error wrapping with eris for methods returning errors
   - Dependency injection setup (fx.Options or dig providers)

3. **File System (`fs.go`)**: Handles directory traversal, module detection, and path computation between API and implementation directories.

4. **Main Entry (`main.go`, `generate.go`)**: CLI interface using urfave/cli that orchestrates the generation process.

### Key Patterns

- **Repository Pattern**: Interfaces must end with "Repository" to be detected
- **Implementation Naming**: Implementations are named with "Impl" suffix (e.g., `Repository` → `repositoryImpl`)
- **File Naming**: Implementation files use snake_case with "_impl.go" suffix
- **Package Naming**: Implementation packages append "impl" to the API package name
- **Dependency Injection**: Supports both Uber's fx and Google's dig frameworks

### Generation Flow

1. Crawl API directory for Go files
2. Parse interfaces ending with "Repository"
3. Check existing implementation files to avoid overwriting
4. Generate new implementations or add missing methods
5. Create/update `repositories.go` with dependency injection setup

### Important Notes

- The tool preserves existing implementations and only adds missing methods
- Generated files include a header comment indicating they are auto-generated
- The tool uses go-tree-sitter for parsing, which is vendored in the project
- When using `--focus`, the stub file (`repositories.go`) is not generated