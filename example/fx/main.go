package main

import (
	"context"
	"fmt"

	"example/api/heisenberg"
	"example/internal"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		internal.Repositories,
		fx.Invoke(func(r heisenberg.ChemistryRepository) {
			if batch, err := r.Cook(context.Background(), heisenberg.Formula{
				Name:   "Blue Crystal",
				Purity: 99.1,
			}); err != nil {
				panic(err)
			} else {
				fmt.Printf("Cooked batch %s with %.1f%% purity\n", batch.ID, batch.Formula.Purity)
			}
		}),
	)
	app.Run()
}

