// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package shrekimpl

import (
	"context"
	"example/api/shrek"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/dig"
)

type SwampDependencies struct {
	dig.In
	// Add dependencies here
}

func NewSwampRepository(deps SwampDependencies) shrek.SwampRepository {
	return &swampRepositoryImpl{
		SwampDependencies: deps,
	}
}

type swampRepositoryImpl struct {
	SwampDependencies
}

func (r *swampRepositoryImpl) CleanSwamp(ctx context.Context) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Swamp.CleanSwamp")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.SwampRepository.CleanSwamp")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.SwampRepository.CleanSwamp")
}

func (r *swampRepositoryImpl) GetSwampStatus(ctx context.Context) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Swamp.GetSwampStatus")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.SwampRepository.GetSwampStatus")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.SwampRepository.GetSwampStatus")
}

func (r *swampRepositoryImpl) EvictIntruders(ctx context.Context, characters []string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Swamp.EvictIntruders")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.SwampRepository.EvictIntruders")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.SwampRepository.EvictIntruders")
}

func (r *swampRepositoryImpl) WelcomeVisitor(ctx context.Context, visitor string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Swamp.WelcomeVisitor")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.SwampRepository.WelcomeVisitor")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.SwampRepository.WelcomeVisitor")
}
