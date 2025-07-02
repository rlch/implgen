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

type Dependencies struct {
	fx.In
	// Add dependencies here
}

var Options = fx.Options(
	fx.Provide(
		NewRepository,
	),
)

func NewRepository(deps Dependencies) waltuh.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}

type repositoryImpl struct {
	Dependencies
}

func (r *repositoryImpl) MakeBreakfast(birthday, kilograms int) waltuh.Waltuh {
	panic("TODO: implement waltuh.Repository.MakeBreakfast")
}

func (r *repositoryImpl) SynthesizeMeth(ctx context.Context, flyPresent, withJesse bool) int {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.SynthesizeMeth")
	defer span.End()
	_ = ctx
	panic("TODO: implement waltuh.Repository.SynthesizeMeth")
}

func (r *repositoryImpl) MakeMoney(ctx context.Context, poundsOfMeth int) (_ int, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.MakeMoney")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.Repository.MakeMoney")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement waltuh.Repository.MakeMoney")
}

func (r *repositoryImpl) DropWaltJrOffAtSchool(ctx context.Context) (_ bool, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.DropWaltJrOffAtSchool")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.Repository.DropWaltJrOffAtSchool")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement waltuh.Repository.DropWaltJrOffAtSchool")
}

func (r *repositoryImpl) KillKrazy8(ctx context.Context, missingPlateShards int) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.KillKrazy8")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.Repository.KillKrazy8")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement waltuh.Repository.KillKrazy8")
}

func (r *repositoryImpl) Get(ctx context.Context, id string) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("waltuh").Start(ctx, "Repository.Get")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "waltuh.Repository.Get")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement waltuh.Repository.Get")
}

func (r *repositoryImpl) Nope() {
	panic("TODO: implement waltuh.Repository.Nope")
}
