// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package waltuhimpl

import (
	"context"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

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
	return "kill confirmed", nil
}
