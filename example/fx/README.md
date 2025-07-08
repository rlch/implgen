# FX Example

This example demonstrates using implgen with Uber's fx dependency injection framework for both Repository and Service patterns.

## Overview

This example project shows how implgen generates implementation files with fx-compatible dependency injection setup. It demonstrates generating implementations for different interface patterns (Repository and Service) in separate directories.

## Project Structure

```
example/fx/
├── api/                           # Interface definitions
│   ├── heisenberg/
│   │   ├── repository.go          # Repository interfaces (ChemistryRepository, MoneyRepository)
│   │   └── service.go             # Service interfaces (NotificationService)
│   └── spongebob/
│       ├── repository.go          # Repository interfaces (JellyfishingRepository, KrustyKrabRepository)
│       └── service.go             # Service interfaces (PattyService, FryService)
├── repository/                    # Generated repository implementations
│   ├── repository.go              # FX dependency injection setup for repositories
│   ├── heisenberg/
│   │   ├── chemistry_impl.go      # ChemistryRepository implementation
│   │   └── money_impl.go          # MoneyRepository implementation
│   └── spongebob/
│       ├── jellyfishing_impl.go   # JellyfishingRepository implementation
│       └── krusty_krab_impl.go    # KrustyKrabRepository implementation
├── service/                       # Generated service implementations
│   ├── service.go                 # FX dependency injection setup for services
│   ├── heisenberg/
│   │   └── notification_impl.go   # NotificationService implementation
│   └── spongebob/
│       ├── patty_impl.go          # PattyService implementation
│       └── fry_impl.go            # FryService implementation
├── gen.go                         # go:generate directive
├── go.mod
├── go.sum
└── main.go                        # Application entry point
```

## Interface Examples

### Repository Pattern

```go
// api/heisenberg/repository.go
type ChemistryRepository interface {
    Cook(ctx context.Context, formula Formula) (*Batch, error)
    GetBatch(ctx context.Context, id string) (*Batch, error)
    OptimizeFormula(formula Formula) (Formula, []string, error)
    TestMethod(ctx context.Context, input string) (string, error)
}

type MoneyRepository interface {
    Launder(ctx context.Context, amount *float64) (*float64, error)
    ProcessPayments(ctx context.Context, amounts []float64) ([]string, error)
}
```

### Service Pattern

```go
// api/heisenberg/service.go
type NotificationService interface {
    SendAlert(ctx context.Context, message string) error
    GetStatus(ctx context.Context, id string) (string, error)
}

// api/spongebob/service.go
type PattyService interface {
    GrillPatty(ctx context.Context, orderID string) error
    IsPattyReady(ctx context.Context, orderID string) (bool, error)
    ServePatty(ctx context.Context, orderID string, customerName string) error
}
```

## Usage

### Generate Repository Implementations

```bash
# Generate all repository implementations in the repository/ directory
./implgen generate --suffix Repository --api api --impl repository

# This generates:
# - repository/repository.go (fx.Options for all repositories)
# - repository/heisenberg/chemistry_impl.go
# - repository/heisenberg/money_impl.go
# - repository/spongebob/jellyfishing_impl.go
# - repository/spongebob/krusty_krab_impl.go
```

### Generate Service Implementations

```bash
# Generate all service implementations in the service/ directory
./implgen generate --suffix Service --api api --impl service

# This generates:
# - service/service.go (fx.Options for all services)
# - service/heisenberg/notification_impl.go
# - service/spongebob/patty_impl.go
# - service/spongebob/fry_impl.go
```

## Generated Features

Each generated implementation includes:

- **Proper package structure** with impl suffix (e.g., `heisenbergimpl`)
- **Dependency injection setup** using fx.In and fx.Options
- **OpenTelemetry tracing** for methods with context.Context
- **Error wrapping** with fault library for methods returning errors
- **Mock generation directives** using moq
- **Stub implementations** with TODO panics for all interface methods

## Key Files

### Repository Stub (`repository/repository.go`)

```go
var Repositories = fx.Options(
    heisenbergimpl.ChemistryOptions,
    heisenbergimpl.MoneyOptions,
    spongebobimpl.JellyfishingOptions,
    spongebobimpl.KrustyKrabOptions,
)
```

### Service Stub (`service/service.go`)

```go
var Services = fx.Options(
    heisenbergimpl.NotificationOptions,
    spongebobimpl.FryOptions,
    spongebobimpl.PattyOptions,
)
```

## Running the Example

```bash
# Build the application
go build -o example .

# Run with repositories and services
./example
```

This example showcases the full flexibility of implgen's generalized interface implementation generation, supporting both Repository and Service patterns in separate, organized directories.
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