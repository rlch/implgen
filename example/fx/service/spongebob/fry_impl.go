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

type FryDependencies struct {
	fx.In
	// Add dependencies here
}

var FryOptions = fx.Options(
	fx.Provide(
		NewFryService,
	),
)

func NewFryService(deps FryDependencies) spongebob.FryService {
	return &fryServiceImpl{
		FryDependencies: deps,
	}
}

type fryServiceImpl struct {
	FryDependencies
}

func (r *fryServiceImpl) CookFries(ctx context.Context, size string) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Fry.CookFries")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.FryService.CookFries"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.FryService.CookFries")
}

func (r *fryServiceImpl) SeasonFries(ctx context.Context, batchID, seasoning string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Fry.SeasonFries")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.FryService.SeasonFries"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.FryService.SeasonFries")
}
