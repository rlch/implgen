package main

import (
	"context"
	"fmt"

	"example/api/waltuh"
	"example/internal"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		internal.Repositories,
		fx.Invoke(func(r waltuh.Repository) {
			if money, err := r.MakeMoney(context.Background(), 100); err != nil {
				panic(err)
			} else {
				fmt.Printf("made $%d money from meth bITCH\n", money)
			}
		}),
	)
	app.Run()
}
