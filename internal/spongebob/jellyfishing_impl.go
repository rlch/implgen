// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package spongebobimpl

import (
	"context"

	"example/api/spongebob"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/fx"
)

type JellyfishingDependencies struct {
	fx.In
	// Add dependencies here
}

var JellyfishingOptions = fx.Options(
	fx.Provide(
		NewJellyfishingRepository,
	),
)

func NewJellyfishingRepository(deps JellyfishingDependencies) spongebob.JellyfishingRepository {
	return &jellyfishingRepositoryImpl{
		JellyfishingDependencies: deps,
	}
}

type jellyfishingRepositoryImpl struct {
	JellyfishingDependencies
}

func (r *jellyfishingRepositoryImpl) CatchJellyfish(ctx context.Context, jellyfish spongebob.Jellyfish) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Jellyfishing.CatchJellyfish")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebob.JellyfishingRepository.CatchJellyfish")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.JellyfishingRepository.CatchJellyfish")
}

func (r *jellyfishingRepositoryImpl) ReleaseJellyfish(ctx context.Context, jellyfish []spongebob.Jellyfish) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Jellyfishing.ReleaseJellyfish")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebob.JellyfishingRepository.ReleaseJellyfish")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.JellyfishingRepository.ReleaseJellyfish")
}

func (r *jellyfishingRepositoryImpl) CountJellyfish(ctx context.Context, area string) (_ int, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Jellyfishing.CountJellyfish")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebob.JellyfishingRepository.CountJellyfish")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.JellyfishingRepository.CountJellyfish")
}

func (r *jellyfishingRepositoryImpl) GetBestSpot(ctx context.Context) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("spongebob").Start(ctx, "Jellyfishing.GetBestSpot")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebob.JellyfishingRepository.GetBestSpot")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement spongebob.JellyfishingRepository.GetBestSpot")
}
