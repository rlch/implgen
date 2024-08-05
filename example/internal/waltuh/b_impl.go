// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package waltuhimpl

import (
	"context"

	"example/api/waltuh"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/fx"
)

type BDependencies struct {
	fx.In
	// Add dependencies here
}

var BOptions = fx.Options(
	fx.Provide(
		NewBRepository,
	),
)

func NewBRepository(deps BDependencies) waltuh.BRepository {
	return &bRepositoryImpl{
		BDependencies: deps,
	}
}

type bRepositoryImpl struct {
	BDependencies
}

func (r *bRepositoryImpl) Yep(ctx context.Context, id string) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "B.Yep")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.BRepository.Yep")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement waltuh.BRepository.Yep")
}

func (r *bRepositoryImpl) Yope() (_ string, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.BRepository.Yope")
		}
	}()
	panic("TODO: implement waltuh.BRepository.Yope")
}
