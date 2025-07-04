package minecraftimpl

import (
	"context"

	"example/api/minecraft"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

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
