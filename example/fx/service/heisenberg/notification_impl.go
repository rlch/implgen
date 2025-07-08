// This file will be automatically regenerated based on the API. Any contract implementations
// will be copied through when generating and new methods will be added to the end.
package heisenbergimpl

import (
	"context"
	"example/api/heisenberg"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/fx"
)

type NotificationDependencies struct {
	fx.In
	// Add dependencies here
}

var NotificationOptions = fx.Options(
	fx.Provide(
		NewNotificationService,
	),
)

func NewNotificationService(deps NotificationDependencies) heisenberg.NotificationService {
	return &notificationServiceImpl{
		NotificationDependencies: deps,
	}
}

type notificationServiceImpl struct {
	NotificationDependencies
}

func (r *notificationServiceImpl) SendAlert(ctx context.Context, message string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Notification.SendAlert")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.NotificationService.SendAlert"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.NotificationService.SendAlert")
}

func (r *notificationServiceImpl) GetStatus(ctx context.Context, id string) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Notification.GetStatus")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.NotificationService.GetStatus"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.NotificationService.GetStatus")
}
