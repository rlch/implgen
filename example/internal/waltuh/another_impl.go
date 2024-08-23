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

type AnotherDependencies struct {
	fx.In
	// Add dependencies here
}

var AnotherOptions = fx.Options(
	fx.Provide(
		NewAnotherRepository,
	),
)

func NewAnotherRepository(deps AnotherDependencies) waltuh.AnotherRepository {
	return &anotherRepositoryImpl{
		AnotherDependencies: deps,
	}
}

type anotherRepositoryImpl struct {
	AnotherDependencies
}

func (r *anotherRepositoryImpl) A() (_ string, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.AnotherRepository.A")
		}
	}()
	panic("TODO: implement waltuh.AnotherRepository.A")
}

func (r *anotherRepositoryImpl) B(ctx context.Context) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Another.B")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.AnotherRepository.B")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement waltuh.AnotherRepository.B")
}

func (r *anotherRepositoryImpl) C() (_ string, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.AnotherRepository.C")
		}
	}()
	panic("TODO: implement waltuh.AnotherRepository.C")
}
