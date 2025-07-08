// This file will be automatically regenerated based on the API. Any contract implementations
// will be copied through when generating and new methods will be added to the end.
package spongebobimpl

import (
	"context"
	"example/api/spongebob"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/fx"
)

type PattyDependencies struct {
	fx.In
	// Add dependencies here
}

var PattyOptions = fx.Options(
	fx.Provide(
		NewPattyService,
	),
)

func NewPattyService(deps PattyDependencies) spongebob.PattyService {
	return &pattyServiceImpl{
		PattyDependencies: deps,
	}
}

type pattyServiceImpl struct {
	PattyDependencies
}

func (r *pattyServiceImpl) GrillPatty(ctx context.Context, orderID string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Patty.GrillPatty")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.PattyService.GrillPatty"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.PattyService.GrillPatty")
}

func (r *pattyServiceImpl) IsPattyReady(ctx context.Context, orderID string) (_ bool, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Patty.IsPattyReady")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.PattyService.IsPattyReady"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.PattyService.IsPattyReady")
}

func (r *pattyServiceImpl) ServePatty(ctx context.Context, orderID, customerName string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Patty.ServePatty")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.PattyService.ServePatty"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.PattyService.ServePatty")
}
