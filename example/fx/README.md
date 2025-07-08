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

### CLI Options

```bash
# Basic generation command
./implgen generate [OPTIONS]

# Available options:
#   --root value       Root directory (default: ".")
#   --api value        API directory relative to root (default: "api")
#   --impl value       Implementation directory relative to root (default: "internal")
#   --suffix value     Interface suffix to detect (default: "Repository")
#   --focus value      Focus on specific packages using glob patterns
#   --dig              Use dig instead of fx for dependency injection
#   --verbose, -v      Enable verbose logging
```

### Generate Repository Implementations

```bash
# Generate all repository implementations in the repository/ directory
./implgen generate --suffix Repository --api api --impl repository

# With verbose output
./implgen generate --suffix Repository --api api --impl repository --verbose

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

# Focus on specific packages only
./implgen generate --suffix Service --api api --impl service --focus "heisenberg/**"

# This generates:
# - service/service.go (fx.Options for all services)
# - service/heisenberg/notification_impl.go
# - service/spongebob/patty_impl.go
# - service/spongebob/fry_impl.go
```

### Using Dig Instead of FX

```bash
# Generate with dig dependency injection
./implgen generate --suffix Repository --api api --impl repository --dig
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

## Recommended Setup with gen.go

Create a `gen.go` file to simplify generation:

```go
package main

//go:generate go run github.com/rlch/implgen generate --suffix Repository --api api --impl repository
//go:generate go run github.com/rlch/implgen generate --suffix Service --api api --impl service

// Alternative: Use local binary (if you have implgen installed)
//go:generate implgen generate --suffix Repository --api api --impl repository
//go:generate implgen generate --suffix Service --api api --impl service

// With additional options:
//go:generate go run github.com/rlch/implgen generate --suffix Repository --api api --impl repository --verbose
//go:generate go run github.com/rlch/implgen generate --suffix Service --api api --impl service --dig
```

Then simply run:

```bash
go generate
```

This will generate both repository and service implementations in one command.

### Advanced Usage Examples

```bash
# Generate only specific packages
./implgen generate --suffix Repository --api api --impl repository --focus "heisenberg/**"

# Use dig instead of fx
./implgen generate --suffix Service --api api --impl service --dig

# Generate Handler implementations (works with any suffix)
./implgen generate --suffix Handler --api api --impl handler

# Specify custom root directory
./implgen generate --suffix Repository --root ./my-project --api interfaces --impl implementations
```