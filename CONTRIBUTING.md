# Contributing to implgen

Thank you for your interest in contributing to implgen! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.23 or later
- C compiler (required for tree-sitter dependencies)
- golangci-lint for code quality checks
- git for version control

### Development Setup

1. **Fork and clone the repository**:
   ```bash
   git clone https://github.com/your-username/implgen.git
   cd implgen
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Build the project**:
   ```bash
   go build -o implgen .
   ```

4. **Run tests**:
   ```bash
   go test ./...
   ```

5. **Run linter**:
   ```bash
   golangci-lint run --config .golangci.yaml ./...
   ```

## Development Workflow

### Code Organization

The codebase is organized into several key files:

- **`main.go`**: CLI entry point and command definitions
- **`generate.go`**: Generation orchestration and workflow
- **`parse.go`**: Interface parsing using tree-sitter
- **`codegen.go`**: Code generation and templating
- **`fs.go`**: File system utilities and path management
- **`example/`**: Working examples for testing and demonstration

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Follow the existing code style and patterns
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   # Run unit tests
   go test ./...
   
   # Test with examples
   cd example/fx
   ../../implgen generate
   go build .
   
   # Run linter
   golangci-lint run --config .golangci.yaml ./...
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and create a pull request**:
   ```bash
   git push origin feature/your-feature-name
   ```

## Code Style Guidelines

### Go Code Style

- Follow standard Go conventions and idioms
- Use `gofumpt` for formatting (stricter than `gofmt`)
- Run `golangci-lint` and address all issues
- Use meaningful variable and function names
- Add comments for public functions and complex logic

### Commit Message Convention

Use conventional commit format:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for adding or updating tests
- `chore:` for maintenance tasks

Examples:
```
feat: add support for generic interfaces
fix: handle nested package imports correctly
docs: update README with new CLI options
refactor: simplify error handling in parser
test: add integration tests for dig example
```

### Error Handling

- Use structured error wrapping with context:
  ```go
  if err != nil {
      return fmt.Errorf("failed to parse file %s: %w", filename, err)
  }
  ```
- Prefer specific error types for different failure modes
- Include relevant context in error messages

### Logging

- Use structured logging with `slog`
- Provide debug-level logging for troubleshooting
- Include relevant context in log messages:
  ```go
  slog.Debug(
      "Parsing repositories",
      slog.String("package", packagePath),
      slog.Int("file_count", len(files)),
  )
  ```

## Testing

### Test Organization

- Unit tests in `*_test.go` files alongside source code
- Integration tests using the `example/` directories
- Test data and fixtures in appropriate subdirectories

### Writing Tests

1. **Unit Tests**: Test individual functions and methods
   ```go
   func TestParseParams(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected Params
       }{
           {
               name:     "simple params",
               input:    "ctx context.Context, id string",
               expected: /* ... */,
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               result := parseParams(tt.input)
               assert.Equal(t, tt.expected, result)
           })
       }
   }
   ```

2. **Integration Tests**: Test complete workflows
   ```go
   func TestGenerateExample(t *testing.T) {
       // Setup test directory
       // Run implgen generate
       // Verify generated files
       // Ensure code compiles
   }
   ```

### Test Guidelines

- Aim for high test coverage of critical paths
- Test both success and failure scenarios
- Use table-driven tests for multiple test cases
- Mock external dependencies where appropriate
- Keep tests fast and deterministic

## Documentation

### Code Documentation

- Add package-level documentation for all packages
- Document public functions, types, and constants
- Use Go doc conventions:
  ```go
  // ParseRepositories extracts Repository interface definitions from Go source code.
  // It uses tree-sitter to parse the syntax tree and extract interface methods,
  // parameters, and return types.
  func ParseRepositories(src []byte, tree *sitter.Tree) ([]*Repository, error) {
      // ...
  }
  ```

### User Documentation

- Update README.md for user-facing changes
- Add examples to EXAMPLES.md for new features
- Update CLI help text for new options
- Keep documentation concise but comprehensive

## Adding New Features

### Supporting New Interface Patterns

To add support for new interface naming patterns:

1. **Update the tree-sitter query** in `parse.go`:
   ```go
   query, queryErr := sitter.NewQuery(language, `
   (type_spec
     name: (type_identifier) @class_name (#match? @class_name "YourPattern$")
     ...
   )`)
   ```

2. **Add tests** for the new pattern
3. **Update documentation** to reflect the new capability

### Adding New Code Generation Features

To add new generated code features:

1. **Extend data structures** in `parse.go`:
   ```go
   type Method struct {
       Ident      string
       Params     Params
       Returns    Params
       NewField   string  // Add your new field
   }
   ```

2. **Update parsing logic** to extract the new information
3. **Modify templates** in `codegen.go` to include the new feature
4. **Update import collection** if new dependencies are needed
5. **Add tests** and update examples

### Supporting New Dependency Injection Frameworks

To add support for a new DI framework:

1. **Add CLI flag** in `main.go`
2. **Update template selection** in `generateRepositoryImpl`
3. **Add new template** for the framework
4. **Update stub generation** in `generateRepositoryStubFile`
5. **Add example** in a new directory
6. **Update documentation**

## Tree-sitter Development

### Understanding Tree-sitter Queries

implgen uses S-expression queries to extract interface definitions:

```scheme
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

### Debugging Tree-sitter

To debug parsing issues:

1. **Enable verbose logging**:
   ```bash
   ./implgen --verbose generate
   ```

2. **Use tree-sitter CLI** to inspect syntax trees:
   ```bash
   tree-sitter parse your-file.go
   ```

3. **Test queries independently**:
   ```bash
   tree-sitter query path/to/grammar.js your-query.scm your-file.go
   ```

### Modifying Queries

When modifying tree-sitter queries:

1. Test with various Go syntax patterns
2. Ensure backward compatibility
3. Add test cases for edge cases
4. Update error handling for new capture groups

## Release Process

### Versioning

- Follow semantic versioning (semver)
- Tag releases with `git tag v1.2.3`
- Update CHANGELOG.md with release notes

### Pre-release Checklist

- [ ] All tests pass
- [ ] Linter passes with no warnings
- [ ] Examples build and run successfully
- [ ] Documentation is up to date
- [ ] Version numbers are updated
- [ ] CHANGELOG.md is updated

## Getting Help

### Community

- GitHub Issues: Report bugs and request features
- GitHub Discussions: Ask questions and share ideas
- Pull Requests: Contribute code and documentation

### Debugging

Common debugging techniques:

1. **Enable verbose logging**:
   ```bash
   ./implgen --verbose generate
   ```

2. **Use a debugger** with your favorite Go debugging tool

3. **Add temporary logging**:
   ```go
   slog.Debug("Debug info", slog.Any("data", yourData))
   ```

4. **Test with minimal examples** to isolate issues

## Architecture Guidelines

### Code Organization Principles

- **Separation of concerns**: Each file has a clear responsibility
- **Single responsibility**: Functions do one thing well
- **Dependency injection**: Minimize global state
- **Error handling**: Consistent error handling patterns
- **Testability**: Design for easy testing

### Performance Considerations

- **Memory efficiency**: Clean up tree-sitter resources
- **Parsing efficiency**: Reuse parsers and queries
- **I/O efficiency**: Minimize file system operations
- **Caching**: Cache expensive operations (module detection)

### Security Considerations

- **Input validation**: Validate all user inputs
- **Path traversal**: Prevent directory traversal attacks
- **Resource limits**: Prevent resource exhaustion
- **Code injection**: Safely handle user-provided templates

## Common Pitfalls

### Tree-sitter Issues

- **Resource leaks**: Always close trees, queries, and cursors
- **Query syntax**: S-expression syntax can be tricky
- **Language binding**: Ensure go-tree-sitter is properly built

### Go Module Issues

- **Import paths**: Ensure correct module paths in generated code
- **Version compatibility**: Test with different Go versions
- **Vendor dependencies**: Handle vendored dependencies correctly

### Template Issues

- **Escaping**: Properly escape template variables
- **Whitespace**: Be careful with template whitespace
- **Error handling**: Handle template execution errors gracefully

## Thank You

Thank you for contributing to implgen! Your contributions help make Go development more productive and enjoyable for everyone.