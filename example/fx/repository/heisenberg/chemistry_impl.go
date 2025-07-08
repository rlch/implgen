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

type ChemistryDependencies struct {
	fx.In
	// Add dependencies here
}

var ChemistryOptions = fx.Options(
	fx.Provide(
		NewChemistryRepository,
	),
)

func NewChemistryRepository(deps ChemistryDependencies) heisenberg.ChemistryRepository {
	return &chemistryRepositoryImpl{
		ChemistryDependencies: deps,
	}
}

type chemistryRepositoryImpl struct {
	ChemistryDependencies
}

func (r *chemistryRepositoryImpl) Cook(ctx context.Context, formula heisenberg.Formula) (_ *heisenberg.Batch, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Chemistry.Cook")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.ChemistryRepository.Cook"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.ChemistryRepository.Cook")
}

func (r *chemistryRepositoryImpl) GetBatch(ctx context.Context, id string) (_ *heisenberg.Batch, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Chemistry.GetBatch")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.ChemistryRepository.GetBatch"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.ChemistryRepository.GetBatch")
}

func (r *chemistryRepositoryImpl) OptimizeFormula(formula heisenberg.Formula) (_ heisenberg.Formula, _ []string, err error) {
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.ChemistryRepository.OptimizeFormula"))
		}
	}()
	panic("TODO: implement heisenberg.ChemistryRepository.OptimizeFormula")
}

func (r *chemistryRepositoryImpl) TestMethod(ctx context.Context, input string) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("heisenberg").Start(ctx, "Chemistry.TestMethod")
	defer func() {
		if err != nil {
			err = fault.Wrap(err, fmsg.With("heisenberg.ChemistryRepository.TestMethod"))
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement heisenberg.ChemistryRepository.TestMethod")
}
