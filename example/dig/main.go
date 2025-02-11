package main

import (
	"context"
	"fmt"

	"example/api/waltuh"
	"example/internal"

	"go.uber.org/dig"
)

func main() {
	dig := dig.New()
	for _, factory := range internal.RepositoryFactories {
		if err := dig.Provide(factory); err != nil {
			panic(err)
		}
	}
	err := dig.Invoke(func(r waltuh.Repository) {
		fmt.Println(r.KillKrazy8(
			context.Background(),
			10,
		))
	})
	if err != nil {
		panic(err)
	}
}
