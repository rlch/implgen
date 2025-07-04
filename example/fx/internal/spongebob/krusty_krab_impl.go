// This file will be automatically regenerated based on the API. Any repository implementations
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

type KrustyKrabDependencies struct {
	fx.In
	// Add dependencies here
}

var KrustyKrabOptions = fx.Options(
	fx.Provide(
		NewKrustyKrabRepository,
	),
)

func NewKrustyKrabRepository(deps KrustyKrabDependencies) spongebob.KrustyKrabRepository {
	return &krustyKrabRepositoryImpl{
		KrustyKrabDependencies: deps,
	}
}

type krustyKrabRepositoryImpl struct {
	KrustyKrabDependencies
}

func (r *krustyKrabRepositoryImpl) CookPatty(ctx context.Context, toppings []string) (_ *spongebob.KrabbyPatty, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "KrustyKrab.CookPatty")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.KrustyKrabRepository.CookPatty"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.KrustyKrabRepository.CookPatty")
}

func (r *krustyKrabRepositoryImpl) ServeCustomer(ctx context.Context, patty spongebob.KrabbyPatty, customerName string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "KrustyKrab.ServeCustomer")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.KrustyKrabRepository.ServeCustomer"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.KrustyKrabRepository.ServeCustomer")
}

func (r *krustyKrabRepositoryImpl) GetMenu(ctx context.Context) (_ []spongebob.MenuItem, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "KrustyKrab.GetMenu")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.KrustyKrabRepository.GetMenu"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.KrustyKrabRepository.GetMenu")
}

func (r *krustyKrabRepositoryImpl) UpdateMenu(ctx context.Context, items []spongebob.MenuItem) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "KrustyKrab.UpdateMenu")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("spongebob.KrustyKrabRepository.UpdateMenu"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.KrustyKrabRepository.UpdateMenu")
}
