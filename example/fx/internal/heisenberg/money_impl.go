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

type MoneyDependencies struct {
	fx.In
	// Add dependencies here
}

var MoneyOptions = fx.Options(
	fx.Provide(
		NewMoneyRepository,
	),
)

func NewMoneyRepository(deps MoneyDependencies) heisenberg.MoneyRepository {
	return &moneyRepositoryImpl{
		MoneyDependencies: deps,
	}
}

type moneyRepositoryImpl struct {
	MoneyDependencies
}

func (r *moneyRepositoryImpl) Launder(ctx context.Context, amount *float64) (_ *float64, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Money.Launder")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.MoneyRepository.Launder"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.MoneyRepository.Launder")
}

func (r *moneyRepositoryImpl) ProcessPayments(ctx context.Context, amounts []float64) (_ []string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Money.ProcessPayments")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.MoneyRepository.ProcessPayments"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.MoneyRepository.ProcessPayments")
}
