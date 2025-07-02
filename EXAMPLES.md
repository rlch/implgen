# Examples

This document provides detailed walkthroughs of the examples included with implgen, showing how to use the tool in real-world scenarios.

## Overview

The `example/` directory contains two complete working examples:

- **`example/fx/`**: Demonstrates using Uber's fx for dependency injection
- **`example/dig/`**: Demonstrates using Google's dig for dependency injection

Both examples showcase the same functionality but with different dependency injection frameworks.

## Example Structure

Each example follows this structure:

```
example/fx/                          # or example/dig/
├── api/                            # Interface definitions
│   ├── spongebob_squarepants/
│   │   └── repository.go
│   └── waltuh/
│       ├── another.go
│       ├── entity.go
│       ├── nested/
│       │   └── repository.go
│       └── repository.go
├── internal/                       # Generated implementations
│   ├── repositories.go             # Dependency injection stub
│   ├── spongebob_squarepants/
│   │   ├── b.go
│   │   └── repository_impl.go
│   └── waltuh/
│       ├── another_impl.go
│       ├── b_impl.go
│       ├── nested/
│       │   └── repository_impl.go
│       └── repository_impl.go
├── gen.go                          # go:generate directive
├── go.mod
├── go.sum
└── main.go                         # Application entry point
```

## Walkthrough: FX Example

Let's walk through the fx example step by step.

### 1. Interface Definitions

#### Simple Repository Interface

```go
// example/fx/api/waltuh/repository.go
package waltuh

import "context"

type Repository interface {
    MakeBreakfast(birthday, kilograms int) Waltuh
    SynthesizeMeth(ctx context.Context, flyPresent bool, withJesse bool) int
    MakeMoney(ctx context.Context, poundsOfMeth int) (int, error)
    DropWaltJrOffAtSchool(ctx context.Context) (bool, error)
    KillKrazy8(ctx context.Context, missingPlateShards int) (string, error)
    Get(ctx context.Context, id string) (string, error)
    Nope()
}
```

#### Multiple Repositories in One Package

```go
// example/fx/api/waltuh/another.go
package waltuh

import "context"

type AnotherRepository interface {
    DoSomething(ctx context.Context) error
}

type BRepository interface {
    Yep(ctx context.Context, id string) (string, error)
    Yope() (string, error)
}
```

#### Nested Package Repository

```go
// example/fx/api/waltuh/nested/repository.go
package nested

import "context"

type Repository interface {
    NestedMethod(ctx context.Context, param string) (string, error)
}
```

### 2. Generation Process

Run implgen in the example directory:

```bash
cd example/fx
../../implgen generate
```

This command:
1. Crawls the `api/` directory for Go files
2. Parses interfaces ending with "Repository"
3. Generates implementation files in `internal/`
4. Creates a `repositories.go` file with dependency injection setup

### 3. Generated Implementation

#### Main Repository Implementation

```go
// example/fx/internal/waltuh/repository_impl.go
package waltuhimpl

import (
    "context"
    "example/api/waltuh"
    "github.com/rotisserie/eris"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/codes"
    "go.uber.org/fx"
)

type Dependencies struct {
    fx.In
    // Add dependencies here
}

var Options = fx.Options(
    fx.Provide(
        NewRepository,
    ),
)

func NewRepository(deps Dependencies) waltuh.Repository {
    return &repositoryImpl{
        Dependencies: deps,
    }
}

type repositoryImpl struct {
    Dependencies
}

func (r *repositoryImpl) MakeBreakfast(birthday, kilograms int) waltuh.Waltuh {
    panic("TODO: implement waltuh.Repository.MakeBreakfast")
}

func (r *repositoryImpl) SynthesizeMeth(ctx context.Context, flyPresent, withJesse bool) int {
    ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.SynthesizeMeth")
    defer span.End()
    _ = ctx
    panic("TODO: implement waltuh.Repository.SynthesizeMeth")
}

func (r *repositoryImpl) MakeMoney(ctx context.Context, poundsOfMeth int) (_ int, err error) {
    ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.MakeMoney")
    defer func() {
        if err != nil {
            err = eris.Wrap(err, "waltuh.Repository.MakeMoney")
            span.SetStatus(codes.Error, "")
            span.RecordError(err)
        }
        span.End()
    }()
    _ = ctx
    panic("TODO: implement waltuh.Repository.MakeMoney")
}

// ... other methods
```

#### Dependency Injection Stub

```go
// example/fx/internal/repositories.go
package internal

//go:generate moq -out=waltuh/nested/mocks.go -pkg=nestedimpl -rm -skip-ensure ../api/waltuh/nested Repository
//go:generate moq -out=spongebob_squarepants/mocks.go -pkg=spongebobsquarepantsimpl -rm -skip-ensure ../api/spongebob_squarepants Repository
//go:generate moq -out=waltuh/mocks.go -pkg=waltuhimpl -rm -skip-ensure ../api/waltuh AnotherRepository BRepository Repository

import (
    spongebobsquarepantsimpl "example/internal/spongebob_squarepants"
    waltuhimpl "example/internal/waltuh"
    nestedimpl "example/internal/waltuh/nested"
    "go.uber.org/fx"
)

var Repositories = fx.Options(
    nestedimpl.Options,
    spongebobsquarepantsimpl.Options,
    waltuhimpl.Options,
    waltuhimpl.AnotherOptions,
    waltuhimpl.BOptions,
)
```

### 4. Integration with FX

The generated code integrates seamlessly with fx:

```go
// example/fx/main.go
package main

import (
    "example/internal"
    "go.uber.org/fx"
)

func main() {
    fx.New(
        internal.Repositories,
        // Add your other modules here
    ).Run()
}
```

## Key Generated Features

### 1. Observability Integration

Methods with `context.Context` parameters automatically get:
- OpenTelemetry span creation
- Automatic error recording and status setting
- Proper span lifecycle management

```go
func (r *repositoryImpl) SomeMethod(ctx context.Context, param string) (_ Result, err error) {
    ctx, span := otel.GetTracerProvider().Tracer("package").Start(ctx, "Repository.SomeMethod")
    defer func() {
        if err != nil {
            err = eris.Wrap(err, "package.Repository.SomeMethod")
            span.SetStatus(codes.Error, "")
            span.RecordError(err)
        }
        span.End()
    }()
    // Implementation here
}
```

### 2. Error Handling

Methods returning `error` get:
- Automatic error wrapping with eris
- Structured error context with package and method names
- Integration with OpenTelemetry error recording

### 3. Dependency Injection

#### FX Pattern
- `Dependencies` struct with `fx.In` tag
- `Options` variable for module registration
- Constructor function following fx conventions

#### Dig Pattern (when using `--dig` flag)
- `Dependencies` struct with `dig.In` tag
- Constructor functions ready for dig container registration
- `RepositoryFactories` slice for batch registration

### 4. Mock Generation

Automatic `go:generate` directives for creating mocks:
- One directive per package
- Includes all repositories in the package
- Uses moq for mock generation
- Proper package naming and output paths

## Dig Example Differences

The dig example (`example/dig/`) generates slightly different code:

### Dependencies Structure

```go
type Dependencies struct {
    dig.In
    // Add dependencies here
}
```

### No Options Variable

Instead of fx.Options, dig example generates a `RepositoryFactories` slice:

```go
var RepositoryFactories = []any{
    NewRepository,
    NewAnotherRepository,
    NewBRepository,
}
```

### Manual Registration

With dig, you manually register the factories:

```go
container := dig.New()
for _, factory := range internal.RepositoryFactories {
    container.Provide(factory)
}
```

## Running the Examples

### Prerequisites

```bash
cd example/fx  # or example/dig
go mod download
```

### Generate Implementations

```bash
# From the example directory
../../implgen generate

# Or with verbose output
../../implgen --verbose generate
```

### Focus on Specific Packages

```bash
# Only generate for waltuh package
../../implgen generate --focus "waltuh"

# Only generate for nested packages
../../implgen generate --focus "**/nested"
```

### Generate Mocks

```bash
go generate ./...
```

### Build and Run

```bash
go build -o example .
./example
```

## Advanced Use Cases

### Custom Directory Structure

If your project uses different directory names:

```bash
../../implgen generate --api services --impl implementations
```

### Incremental Updates

When you add new methods to existing interfaces:

1. Update the interface in the API directory
2. Run `implgen generate`
3. Only new methods are added to existing implementation files
4. Existing method implementations are preserved

### Generic Types

implgen handles generic types in interfaces:

```go
type Repository[T any] interface {
    Create(ctx context.Context, item T) error
    Get(ctx context.Context, id string) (T, error)
}
```

Generated implementation:

```go
type repositoryImpl[T any] struct {
    Dependencies[T]
}

func (r *repositoryImpl[T]) Create(ctx context.Context, item T) error {
    // Implementation
}
```

## Best Practices

### 1. Interface Design

- Keep interfaces focused and cohesive
- Use context.Context for cancellation and tracing
- Return errors for operations that can fail
- Follow Go naming conventions

### 2. Implementation Organization

- Add your dependencies to the `Dependencies` struct
- Implement methods one at a time, replacing `panic` statements
- Use the generated error wrapping and tracing
- Keep business logic separate from repository logic

### 3. Testing

- Use the generated mocks for unit testing
- Test repository implementations with integration tests
- Leverage the observability features for debugging

### 4. Maintenance

- Re-run implgen when interface definitions change
- Version control both API and implementation files
- Use `--focus` for large projects to avoid unnecessary regeneration

## Troubleshooting

### Interface Not Detected

Check that:
- Interface name ends with "Repository"
- File is in the correct API directory
- Go syntax is valid

### Import Path Issues

Ensure:
- `go.mod` is properly configured
- Module path matches your project structure
- Run `go mod tidy` after generation

### Build Errors

Common solutions:
- Install required dependencies: `go mod download`
- Ensure C compiler is available for tree-sitter
- Check that generated imports are correct