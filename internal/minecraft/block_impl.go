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

type BlockDependencies struct {
	dig.In
	// Add dependencies here
}

func NewBlockRepository(deps BlockDependencies) minecraft.BlockRepository {
	return &blockRepositoryImpl{
		BlockDependencies: deps,
	}
}

type blockRepositoryImpl struct {
	BlockDependencies
}

func (r *blockRepositoryImpl) PlaceBlock(ctx context.Context, blockType minecraft.BlockType, pos minecraft.Coordinates) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Block.PlaceBlock")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.BlockRepository.PlaceBlock")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.BlockRepository.PlaceBlock")
}

func (r *blockRepositoryImpl) GetBlock(ctx context.Context, pos minecraft.Coordinates) (_ *minecraft.Block, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Block.GetBlock")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.BlockRepository.GetBlock")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.BlockRepository.GetBlock")
}

func (r *blockRepositoryImpl) GetArea(ctx context.Context, from, to minecraft.Coordinates) (_ []minecraft.Block, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Block.GetArea")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.BlockRepository.GetArea")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.BlockRepository.GetArea")
}

func (r *blockRepositoryImpl) BulkPlace(ctx context.Context, blocks map[minecraft.Coordinates]minecraft.BlockType) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Block.BulkPlace")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.BlockRepository.BulkPlace")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.BlockRepository.BulkPlace")
}
