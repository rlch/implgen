// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package commonimpl

import (
	"context"
	"example/api/common"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/dig"
)

type BaseDependencies struct {
	dig.In
	// Add dependencies here
}

func NewBaseRepository(deps BaseDependencies) common.BaseRepository {
	return &baseRepositoryImpl{
		BaseDependencies: deps,
	}
}

type baseRepositoryImpl struct {
	BaseDependencies
}

func (r *baseRepositoryImpl) GetMetrics(ctx context.Context) (_ map[string]int64, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("common").Start(ctx, "Base.GetMetrics")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "common.BaseRepository.GetMetrics")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement common.BaseRepository.GetMetrics")
}
