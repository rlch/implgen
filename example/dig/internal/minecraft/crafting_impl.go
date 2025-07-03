// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package minecraftimpl

import (
	"context"
	"example/api/minecraft"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/dig"
)

type CraftingDependencies struct {
	dig.In
	// Add dependencies here
}

func NewCraftingRepository(deps CraftingDependencies) minecraft.CraftingRepository {
	return &craftingRepositoryImpl{
		CraftingDependencies: deps,
	}
}

type craftingRepositoryImpl struct {
	CraftingDependencies
}

func (r *craftingRepositoryImpl) Craft(ctx context.Context, recipe minecraft.CraftingRecipe, inventory map[minecraft.ItemID]int) (_ map[minecraft.ItemID]int, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Crafting.Craft")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.CraftingRepository.Craft")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.CraftingRepository.Craft")
}

func (r *craftingRepositoryImpl) TransferItems(ctx context.Context, fromPlayer, toPlayer string, items map[minecraft.ItemID]int) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Crafting.TransferItems")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.CraftingRepository.TransferItems")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.CraftingRepository.TransferItems")
}

func (r *craftingRepositoryImpl) GetRecipes(ctx context.Context) (_ []minecraft.CraftingRecipe, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Crafting.GetRecipes")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.CraftingRepository.GetRecipes")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.CraftingRepository.GetRecipes")
}
