# FX Example

This example demonstrates using implgen with Uber's fx dependency injection framework.

## Overview

This example project shows how implgen generates implementation files with fx-compatible dependency injection setup. It includes multiple repository interfaces in different packages to showcase various patterns.

## Project Structure

```
example/fx/
├── api/                           # Interface definitions
│   ├── spongebob_squarepants/
│   │   └── repository.go          # Simple repository interface
│   └── waltuh/
│       ├── another.go             # Multiple repositories in one package
│       ├── entity.go              # Entity definitions
│       ├── nested/
│       │   └── repository.go      # Nested package repository
│       └── repository.go          # Main repository with various method patterns
├── internal/                      # Generated implementations
│   ├── repositories.go            # FX dependency injection setup
│   ├── spongebob_squarepants/
│   │   ├── b.go                   # Additional implementation
│   │   └── repository_impl.go     # Generated implementation
│   └── waltuh/
│       ├── another_impl.go        # Generated implementation
│       ├── b_impl.go              # Generated implementation
│       ├── nested/
│       │   └── repository_impl.go # Nested implementation
│       └── repository_impl.go     # Main repository implementation
├── gen.go                         # go:generate directive
├── go.mod
├── go.sum
└── main.go                        # Application entry point
```

## Interface Examples

### Simple Repository

```go
// api/spongebob_squarepants/repository.go
type Repository interface {
    GetKrabbyPatty(ctx context.Context) (KrabbyPatty, error)
    MakeKrabbyPatty(ctx context.Context, patty KrabbyPatty) error
}
```

### Multiple Repositories

```go
// api/waltuh/another.go
type AnotherRepository interface {
    DoSomething(ctx context.Context) error
}

type BRepository interface {
    Yep(ctx context.Context, id string) (string, error)
    Yope() (string, error)
}
```

### Complex Repository with Various Patterns

```go
// api/waltuh/repository.go
type Repository interface {
    // Method without context (no tracing)
    MakeBreakfast(birthday, kilograms int) Waltuh
    
    // Method with context (gets tracing)
    SynthesizeMeth(ctx context.Context, flyPresent bool, withJesse bool) int
    
    // Method with context and error (gets tracing + error wrapping)
    MakeMoney(ctx context.Context, poundsOfMeth int) (int, error)
    
    // Method with no parameters or return values
    Nope()
}
```

## Generated Code Features

### FX Dependency Injection

Each repository gets:

1. **Dependencies struct with fx.In**:
   ```go
   type Dependencies struct {
       fx.In
       // Add dependencies here
   }
   ```

2. **fx.Options export**:
   ```go
   var Options = fx.Options(
       fx.Provide(
           NewRepository,
       ),
   )
   ```

3. **Constructor function**:
   ```go
   func NewRepository(deps Dependencies) waltuh.Repository {
       return &repositoryImpl{
           Dependencies: deps,
       }
   }
   ```

### Implementation Structure

```go
type repositoryImpl struct {
    Dependencies
}
```

### Method Implementations

#### Without Context
```go
func (r *repositoryImpl) MakeBreakfast(birthday, kilograms int) waltuh.Waltuh {
    panic("TODO: implement waltuh.Repository.MakeBreakfast")
}
```

#### With Context (Gets OpenTelemetry Tracing)
```go
func (r *repositoryImpl) SynthesizeMeth(ctx context.Context, flyPresent, withJesse bool) int {
    ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.SynthesizeMeth")
    defer span.End()
    _ = ctx
    panic("TODO: implement waltuh.Repository.SynthesizeMeth")
}
```

#### With Context and Error (Gets Tracing + Error Wrapping)
```go
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
```

### Central FX Module

The `internal/repositories.go` file provides a central fx module:

```go
var Repositories = fx.Options(
    nestedimpl.Options,
    spongebobsquarepantsimpl.Options,
    waltuhimpl.Options,
    waltuhimpl.AnotherOptions,
    waltuhimpl.BOptions,
)
```

### Mock Generation

The file also includes go:generate directives for mock creation:

```go
//go:generate moq -out=waltuh/nested/mocks.go -pkg=nestedimpl -rm -skip-ensure ../api/waltuh/nested Repository
//go:generate moq -out=spongebob_squarepants/mocks.go -pkg=spongebobsquarepantsimpl -rm -skip-ensure ../api/spongebob_squarepants Repository
//go:generate moq -out=waltuh/mocks.go -pkg=waltuhimpl -rm -skip-ensure ../api/waltuh AnotherRepository BRepository Repository
```

## Running the Example

### Generate Implementations

From the example/fx directory:

```bash
# Generate all implementations
../../implgen generate

# Or with verbose output
../../implgen --verbose generate

# Focus on specific packages
../../implgen generate --focus "waltuh/**"
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

## Integration with FX Application

```go
package main

import (
    "context"
    "example/internal"
    "go.uber.org/fx"
)

func main() {
    app := fx.New(
        // Include all generated repositories
        internal.Repositories,
        
        // Add your business logic modules
        fx.Provide(
            NewBusinessLogic,
            NewHTTPServer,
        ),
        
        // Add lifecycle hooks
        fx.Invoke(func(lc fx.Lifecycle, server *HTTPServer) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    return server.Start()
                },
                OnStop: func(ctx context.Context) error {
                    return server.Stop()
                },
            })
        }),
    )
    
    app.Run()
}
```

## Testing with Generated Mocks

```go
func TestBusinessLogic(t *testing.T) {
    mockRepo := &waltuhimpl.RepositoryMock{
        MakeMoneyFunc: func(ctx context.Context, poundsOfMeth int) (int, error) {
            return poundsOfMeth * 1000, nil
        },
    }
    
    logic := NewBusinessLogic(mockRepo)
    result, err := logic.ProcessMeth(context.Background(), 5)
    
    assert.NoError(t, err)
    assert.Equal(t, 5000, result)
    assert.Len(t, mockRepo.MakeMoneyCalls(), 1)
}
```

## Key Benefits

1. **Zero Boilerplate**: No manual dependency injection setup
2. **Observability**: Automatic tracing for context-aware methods
3. **Error Handling**: Structured error wrapping with context
4. **Type Safety**: Full compile-time type checking
5. **Testing**: Generated mocks for easy unit testing
6. **Incremental**: Preserves existing implementations when adding new methods