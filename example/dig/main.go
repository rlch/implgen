package main

import (
	"context"
	"fmt"

	"example/api/shrek"
	"example/internal"

	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	// Register all repository factories
	for _, factory := range internal.RepositoryFactories {
		if err := container.Provide(factory); err != nil {
			panic(err)
		}
	}

	// Use the repository
	if err := container.Invoke(func(r shrek.SwampRepository) {
		if err := r.CleanSwamp(context.Background()); err != nil {
			panic(err)
		} else {
			fmt.Println("Cleaned the swamp successfully!")
		}
	}); err != nil {
		panic(err)
	}
}