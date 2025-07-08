# implgen

A flexible Go code generator that automatically creates implementation boilerplate for Go interfaces with configurable patterns.

## Overview

`implgen` streamlines the development of Go applications by automatically generating implementation files for interfaces. It parses interface definitions with configurable suffixes (Repository, Service, Handler, etc.) from an API directory and creates corresponding implementation files with proper dependency injection, observability, and error handling.

## Features

- 🎯 **Configurable Patterns**: Support for any interface suffix (Repository, Service, Handler, etc.)
- 🏗️ **Dependency Injection**: Supports both fx and dig frameworks
- 📊 **Observability**: Automatic OpenTelemetry tracing for context methods
- 🔄 **Incremental Updates**: Preserves existing implementations, adds missing methods

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
│       ├── repository.go      # Repository interfaces
│       └── service.go         # Service interfaces
├── repository/                # Repository implementations
│   ├── repository.go         # fx.Options for repositories
│   └── user/
│       └── user_repository_impl.go
├── service/                   # Service implementations
│   ├── service.go            # fx.Options for services
│   └── user/
│       └── user_service_impl.go
└── main.go
```

### 2. Define Your Interfaces

Create interfaces with configurable suffixes in your API directory:

```go
// api/user/repository.go
package user

import "context"

type UserRepository interface {
    Create(ctx context.Context, user User) error
    GetByID(ctx context.Context, id string) (User, error)
    Update(ctx context.Context, user User) error
    Delete(ctx context.Context, id string) error
}
```

```go
// api/user/service.go
package user

import "context"

type UserService interface {
    Register(ctx context.Context, user User) error
    Authenticate(ctx context.Context, email, password string) (User, error)
    SendWelcomeEmail(ctx context.Context, userID string) error
}
```

### 3. Generate Implementations

```bash
# Generate repository implementations
./implgen generate --suffix Repository --api api --impl repository

# Generate service implementations  
./implgen generate --suffix Service --api api --impl service
```

This creates separate implementation directories with:

- Proper package structure and imports
- Dependency injection setup
- Method stubs with TODO comments
- OpenTelemetry tracing for context-aware methods
- Error wrapping for methods returning errors

### 4. Generated Output

```go
// repository/user/user_repository_impl.go
package userimpl

import (
    "context"
    "example/api/user"
    "github.com/Southclaws/fault"
    "github.com/Southclaws/fault/fmsg"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/codes"
    "go.uber.org/fx"
)

type UserDependencies struct {
    fx.In
    // Add dependencies here
}

var UserOptions = fx.Options(fx.Provide(NewUserRepository))

func NewUserRepository(deps UserDependencies) user.UserRepository {
    return &userRepositoryImpl{UserDependencies: deps}
}

type userRepositoryImpl struct {
    UserDependencies
}

func (r *userRepositoryImpl) Create(ctx context.Context, user user.User) (err error) {
    ctx, span := otel.GetTracerProvider().Tracer("user").Start(ctx, "User.Create")
    defer func() {
        if err != nil {
            err = fault.Wrap(err, fmsg.With("user.UserRepository.Create"))
            span.SetStatus(codes.Error, "")
            span.RecordError(err)
        }
        span.End()
    }()
    panic("TODO: implement user.UserRepository.Create")
}
```

### 5. Recommended Setup with gen.go

Create a `gen.go` file for easy generation:

```go
package main

//go:generate go run github.com/rlch/implgen generate --suffix Repository --api api --impl repository
//go:generate go run github.com/rlch/implgen generate --suffix Service --api api --impl service
```

Then run:

```bash
go generate
```

## CLI Usage

### Basic Commands

```bash
# Generate repository implementations (default behavior)
./implgen generate

# Generate with specific suffix and directory
./implgen generate --suffix Repository --api api --impl repository

# Generate service implementations
./implgen generate --suffix Service --api api --impl service

# Generate handler implementations
./implgen generate --suffix Handler --api api --impl handler

# Use custom directories
./implgen generate --suffix Repository --root . --api services --impl implementations

# Focus on specific packages
./implgen generate --suffix Repository --impl repository --focus "user/**"
./implgen generate --suffix Service --impl service --focus "user/**" --focus "order/**"

# Use dig instead of fx for dependency injection
./implgen generate --suffix Repository --impl repository --dig

# Enable verbose logging
./implgen --verbose generate --suffix Service --impl service
```

### Command Reference

| Flag        | Description                                    | Default        |
| ----------- | ---------------------------------------------- | -------------- |
| `--root`    | Root directory for the project                 | `.`            |
| `--api`     | API directory relative to root                 | `api`          |
| `--impl`    | Implementation directory relative to root      | `internal`     |
| `--suffix`  | Interface suffix to detect (Repository, Service, etc.) | `Repository` |
| `--focus`   | Target specific packages with glob patterns    | (all packages) |
| `--dig`     | Use dig instead of fx for dependency injection | `false`        |
| `--verbose` | Enable verbose logging                         | `false`        |

## Examples

See the `example/fx/` directory for a complete working example demonstrating:

- Repository and Service patterns in separate directories
- Multiple interfaces per package
- Proper fx dependency injection setup
- Generated implementations with tracing and error handling

## Troubleshooting

### Common Issues

**Interface not detected**

- Ensure interface name ends with the specified suffix (default: "Repository")
- Check that the file is in the correct API directory
- Verify Go syntax is valid
- Use `--verbose` to see what files are being processed

**Build errors with tree-sitter**

- Ensure you have a C compiler installed
- Use `GOFLAGS=-mod=mod` when building (required for tree-sitter)
- Tree-sitter dependencies are vendored in the project

**Generated code has import issues**

- Run `go mod tidy` after generation
- Ensure your module path is correctly set in go.mod
- Check that API interfaces are in proper Go packages

### Debug Mode

Enable verbose logging to see what implgen is doing:

```bash
./implgen --verbose generate --suffix Repository --impl repository
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
