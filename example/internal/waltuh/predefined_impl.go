package waltuhimpl

import (
	"context"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

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
	return false, nil
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
	return "done", nil
}
