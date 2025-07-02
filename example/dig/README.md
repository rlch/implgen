# Dig Example

This example demonstrates using implgen with Google's dig dependency injection framework.

## Overview

This example project shows how implgen generates implementation files with dig-compatible dependency injection setup. It mirrors the fx example but uses dig's container-based approach instead of fx's declarative module system.

## Project Structure

```
example/dig/
├── api/                           # Interface definitions (same as fx example)
│   ├── generic/
│   │   └── repository.go          # Generic repository interface
│   └── waltuh/
│       ├── another.go             # Multiple repositories
│       ├── entity.go              # Entity definitions
│       ├── nested/
│       │   └── repository.go      # Nested package repository
│       └── repository.go          # Main repository
├── internal/                      # Generated implementations
│   ├── repositories.go            # Dig factory functions setup
│   ├── generic/
│   │   ├── mocks.go               # Generated mocks
│   │   ├── multigenerics_impl.go  # Multi-generic implementation
│   │   ├── nogenerics_impl.go     # Non-generic implementation
│   │   └── repository_impl.go     # Main generic implementation
│   └── waltuh/
│       ├── another_impl.go        # Generated implementation
│       ├── b_impl.go              # Generated implementation
│       ├── existing_impl.go       # Pre-existing implementation
│       ├── mocks.go               # Generated mocks
│       ├── nested/
│       │   ├── mocks.go           # Nested mocks
│       │   └── repository_impl.go # Nested implementation
│       └── repository_impl.go     # Main repository implementation
├── gen.go                         # go:generate directive
├── go.mod
├── go.sum
└── main.go                        # Application entry point with dig container
```

## Key Differences from FX Example

### 1. Dependencies Structure

Instead of fx.In, uses dig.In:

```go
type Dependencies struct {
    dig.In
    // Add dependencies here
}
```

### 2. No Options Export

Dig doesn't use module-based registration, so no `Options` variable is generated.

### 3. Factory Functions Collection

Instead of fx.Options, a slice of factory functions is provided:

```go
var RepositoryFactories = []any{
    NewRepository,
    NewAnotherRepository,
    NewBRepository,
}
```

### 4. Manual Container Registration

Unlike fx's automatic module registration, dig requires manual registration:

```go
func main() {
    container := dig.New()
    
    // Register all repository factories
    for _, factory := range internal.RepositoryFactories {
        if err := container.Provide(factory); err != nil {
            log.Fatal(err)
        }
    }
    
    // Register other dependencies
    container.Provide(NewBusinessLogic)
    container.Provide(NewHTTPServer)
    
    // Invoke your application
    if err := container.Invoke(func(server *HTTPServer) {
        server.Start()
    }); err != nil {
        log.Fatal(err)
    }
}
```

## Generated Code Features

### Dig Dependency Injection

Each repository gets:

1. **Dependencies struct with dig.In**:
   ```go
   type Dependencies struct {
       dig.In
       // Add dependencies here
   }
   ```

2. **Constructor function** (same as fx):
   ```go
   func NewRepository(deps Dependencies) waltuh.Repository {
       return &repositoryImpl{
           Dependencies: deps,
       }
   }
   ```

3. **No Options variable** (dig uses direct registration)

### Central Factory Collection

The `internal/repositories.go` file provides a collection of factory functions:

```go
var RepositoryFactories = []any{
    NewRepository,
    NewAnotherRepository,
    NewBRepository,
    // Generic repositories are excluded as they need explicit type parameters
}
```

**Note**: Generic repositories are excluded from the factory slice because dig requires explicit type instantiation.

### Generic Repository Handling

For generic repositories, manual registration is required:

```go
// Generic repository must be registered with specific types
container.Provide(func(deps Dependencies) generic.Repository[string] {
    return NewRepository[string](deps)
})

container.Provide(func(deps Dependencies) generic.Repository[int] {
    return NewRepository[int](deps)
})
```

## Running the Example

### Generate Implementations

From the example/dig directory:

```bash
# Generate implementations with dig flag
../../implgen generate --dig

# Or with verbose output
../../implgen --verbose generate --dig
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

## Integration with Dig Container

```go
package main

import (
    "log"
    "example/internal"
    "go.uber.org/dig"
)

func main() {
    container := dig.New()
    
    // Register all repository factories
    for _, factory := range internal.RepositoryFactories {
        if err := container.Provide(factory); err != nil {
            log.Fatalf("Failed to register factory: %v", err)
        }
    }
    
    // Register generic repositories manually
    if err := container.Provide(func(deps Dependencies) generic.Repository[string] {
        return NewRepository[string](deps)
    }); err != nil {
        log.Fatalf("Failed to register generic repository: %v", err)
    }
    
    // Register business logic
    if err := container.Provide(NewBusinessLogic); err != nil {
        log.Fatalf("Failed to register business logic: %v", err)
    }
    
    // Start the application
    if err := container.Invoke(func(
        repo waltuh.Repository,
        anotherRepo waltuh.AnotherRepository,
        logic *BusinessLogic,
    ) {
        // Use your repositories and business logic
        logic.Run()
    }); err != nil {
        log.Fatalf("Failed to start application: %v", err)
    }
}
```

## Testing with Dig

Dig makes testing straightforward with scoped containers:

```go
func TestBusinessLogic(t *testing.T) {
    container := dig.New()
    
    // Provide mock implementation
    container.Provide(func() waltuh.Repository {
        return &waltuhimpl.RepositoryMock{
            MakeMoneyFunc: func(ctx context.Context, poundsOfMeth int) (int, error) {
                return poundsOfMeth * 1000, nil
            },
        }
    })
    
    // Provide business logic
    container.Provide(NewBusinessLogic)
    
    // Test the business logic
    container.Invoke(func(logic *BusinessLogic) {
        result, err := logic.ProcessMeth(context.Background(), 5)
        assert.NoError(t, err)
        assert.Equal(t, 5000, result)
    })
}
```

## Advanced Dig Features

### Named Dependencies

Dig supports named dependencies for multiple implementations:

```go
// Register multiple implementations of the same interface
container.Provide(NewPrimaryRepository, dig.Name("primary"))
container.Provide(NewSecondaryRepository, dig.Name("secondary"))

// Inject specific implementations
type BusinessLogic struct {
    Primary   waltuh.Repository `dig:"name:'primary'"`
    Secondary waltuh.Repository `dig:"name:'secondary'"`
}
```

### Groups

Dig supports grouping dependencies:

```go
// Register repositories as a group
container.Provide(NewRepository, dig.Group("repositories"))
container.Provide(NewAnotherRepository, dig.Group("repositories"))

// Inject all repositories in the group
type AllRepositories struct {
    Repositories []waltuh.Repository `dig:"group:'repositories'"`
}
```

## Dig vs FX Comparison

| Feature | FX | Dig |
|---------|-----|-----|
| **API Style** | Declarative modules | Imperative container |
| **Registration** | Automatic via Options | Manual via Provide calls |
| **Error Handling** | Startup validation | Runtime validation |
| **Lifecycle** | Built-in lifecycle hooks | Manual lifecycle management |
| **Testing** | Module-based test setup | Container-based test setup |
| **Performance** | Slightly higher overhead | Lower overhead |
| **Learning Curve** | Higher (more concepts) | Lower (simpler API) |

## When to Use Dig

Choose dig when you:

- Prefer imperative dependency registration
- Need fine-grained control over container setup
- Want minimal framework overhead
- Are building simpler applications
- Need custom dependency resolution logic

## Key Benefits

1. **Explicit Control**: Manual registration gives you full control
2. **Performance**: Lower overhead than fx
3. **Simplicity**: Fewer concepts to learn
4. **Flexibility**: Easy to customize dependency resolution
5. **Testing**: Simple container-based testing setup
6. **Observability**: Same automatic tracing and error handling as fx example