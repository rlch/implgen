# implgen

A Go code generator that automatically creates implementation boilerplate for Go interfaces following the Repository pattern.

## Overview

`implgen` streamlines the development of Go applications by automatically generating implementation files for interfaces. It parses interface definitions (ending with "Repository") from an API directory and creates corresponding implementation files with proper dependency injection, observability, and error handling.

## Features

- 🔍 **Smart Parsing**: Uses tree-sitter to accurately parse Go interfaces
- 🏗️ **Dependency Injection**: Supports both Uber's fx and Google's dig frameworks
- 📊 **Observability**: Automatically adds OpenTelemetry tracing for methods with context
- 🚨 **Error Handling**: Integrates eris for structured error wrapping
- 🔄 **Incremental Updates**: Preserves existing implementations while adding missing methods
- 🎯 **Focused Generation**: Target specific packages with glob patterns

## Installation

### Prerequisites

- Go 1.23 or later

### Build from Source

```bash
git clone https://github.com/rlch/implgen
cd implgen
go build -o implgen .
```

## Quick Start

### 1. Project Structure

Organize your Go project with API interfaces and implementation directories:

```
your-project/
├── api/
│   └── user/
│       └── repository.go      # Interface definitions
├── internal/
│   └── user/
│       └── repository_impl.go # Generated implementations
└── main.go
```

### 2. Define Your Interface

Create an interface ending with "Repository" in your API directory:

```go
// api/user/repository.go
package user

import "context"

type Repository interface {
    Create(ctx context.Context, user User) error
    GetByID(ctx context.Context, id string) (User, error)
    Update(ctx context.Context, user User) error
    Delete(ctx context.Context, id string) error
}
```

### 3. Generate Implementation

```bash
./implgen generate
```

This creates `internal/user/repository_impl.go` with:

- Proper package structure and imports
- Dependency injection setup
- Method stubs with TODO comments
- OpenTelemetry tracing for context-aware methods
- Error wrapping for methods returning errors

### 4. Generated Output

```go
// internal/user/repository_impl.go
package userimpl

import (
    "context"
    "example/api/user"
    "github.com/rotisserie/eris"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/codes"
    "go.uber.org/fx"
)

type Dependencies struct {
    fx.In
    // Add dependencies here
}

var Options = fx.Options(fx.Provide(NewRepository))

func NewRepository(deps Dependencies) user.Repository {
    return &repositoryImpl{Dependencies: deps}
}

type repositoryImpl struct {
    Dependencies
}

func (r *repositoryImpl) Create(ctx context.Context, user user.User) (err error) {
    ctx, span := otel.GetTracerProvider().Tracer("user").Start(ctx, "Repository.Create")
    defer func() {
        if err != nil {
            err = eris.Wrap(err, "user.Repository.Create")
            span.SetStatus(codes.Error, "")
            span.RecordError(err)
        }
        span.End()
    }()
    panic("TODO: implement user.Repository.Create")
}
```

## CLI Usage

### Basic Commands

```bash
# Generate all implementations
./implgen generate

# Use custom directories
./implgen generate --root . --api services --impl implementations

# Focus on specific packages
./implgen generate --focus "user/**"
./implgen generate --focus "user/**" --focus "order/**"

# Use dig instead of fx for dependency injection
./implgen generate --dig

# Enable verbose logging
./implgen --verbose generate
```

### Command Reference

| Flag        | Description                                    | Default        |
| ----------- | ---------------------------------------------- | -------------- |
| `--root`    | Root directory for the project                 | `.`            |
| `--api`     | API directory relative to root                 | `api`          |
| `--impl`    | Implementation directory relative to root      | `internal`     |
| `--focus`   | Target specific packages with glob patterns    | (all packages) |
| `--dig`     | Use dig instead of fx for dependency injection | `false`        |
| `--verbose` | Enable verbose logging                         | `false`        |

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   API Layer     │    │   Parser Engine  │    │ Code Generator  │
│                 │───▶│                  │───▶│                 │
│ Interface Defs  │    │  Tree-sitter     │    │ Template Engine │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
                                                         ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Dependency      │◀───│  File System     │◀───│ Implementation  │
│ Injection       │    │                  │    │     Layer       │
│ (fx/dig)        │    │ Path Management  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Core Components

1. **Parser**: Extracts interface definitions using go-tree-sitter
2. **Generator**: Creates implementation files with proper structure
3. **File System**: Manages paths and module detection
4. **CLI**: Orchestrates the generation process

### Key Conventions

- **Interface Naming**: Must end with "Repository"
- **Implementation Naming**: `repositoryImpl` (lowercase first letter + "Impl")
- **File Naming**: `snake_case_impl.go`
- **Package Naming**: API package name + "impl" suffix

## Examples

See the `example/` directory for complete working examples:

- `example/fx/`: Using Uber's fx for dependency injection
- `example/dig/`: Using Google's dig for dependency injection

## Troubleshooting

### Common Issues

**Interface not detected**

- Ensure interface name ends with "Repository"
- Check that the file is in the correct API directory
- Verify Go syntax is valid

**Build errors with tree-sitter**

- Ensure you have a C compiler installed
- Tree-sitter dependencies are vendored in the project

**Generated code has import issues**

- Run `go mod tidy` after generation
- Ensure your module path is correctly set in go.mod

### Debug Mode

Enable verbose logging to see what implgen is doing:

```bash
./implgen --verbose generate
```

## Development

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run --config .golangci.yaml ./...
```

### Building

```bash
go build -o implgen .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass and code is linted
6. Submit a pull request
   j

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [ ] Generate test boilerplate
- [ ] Support for additional interface patterns
- [ ] Integration with popular Go frameworks
- [ ] IDE plugins for seamless development
