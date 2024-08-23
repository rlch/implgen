package waltuhimpl_test

import (
	"context"
	"testing"

	"example/api/waltuh"
	waltuhimpl "example/internal/waltuh"

	fxutils "github.com/MathGaps/core/pkg/utils/fx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestRepository(t *testing.T) {
	t.Parallel()
	type Deps struct {
		fx.In
		Repository waltuh.Repository
	}
	ctx := context.Background()
	_ = ctx
	app := fxtest.New(t,
		fx.Provide(
		// Add dependencies here
		),
		waltuhimpl.Options,

		fxutils.Test(t, func(t *testing.T, d Deps) {
			_ = d
			// Initialization code here
			t.Run("MakeBreakfast", func(t *testing.T) {
				t.Skip()
			})

			t.Run("SynthesizeMeth", func(t *testing.T) {
				t.Skip()
			})
		}),
	)
	app.Run()
}
