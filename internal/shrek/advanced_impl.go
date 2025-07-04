// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package shrekimpl

import (
	"context"
	"io"

	"example/api/shrek"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/dig"
)

type AdvancedDependencies struct {
	dig.In
	// Add dependencies here
}

func NewAdvancedRepository(deps AdvancedDependencies) shrek.AdvancedRepository {
	return &advancedRepositoryImpl{
		AdvancedDependencies: deps,
	}
}

type advancedRepositoryImpl struct {
	AdvancedDependencies
}

func (r *advancedRepositoryImpl) FormParty(ctx context.Context, leader string, members ...string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Advanced.FormParty")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.AdvancedRepository.FormParty")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.AdvancedRepository.FormParty")
}

func (r *advancedRepositoryImpl) ProcessCharacters(ctx context.Context, filter func(string) bool) (_ []string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Advanced.ProcessCharacters")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.AdvancedRepository.ProcessCharacters")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.AdvancedRepository.ProcessCharacters")
}

func (r *advancedRepositoryImpl) ExportStory(ctx context.Context, questID string, writer io.Writer) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Advanced.ExportStory")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.AdvancedRepository.ExportStory")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.AdvancedRepository.ExportStory")
}

func (r *advancedRepositoryImpl) ImportLegend(ctx context.Context, reader io.Reader) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Advanced.ImportLegend")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.AdvancedRepository.ImportLegend")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.AdvancedRepository.ImportLegend")
}

func (r *advancedRepositoryImpl) GetMagicMirrorWisdom() string {
	panic("TODO: implement shrek.AdvancedRepository.GetMagicMirrorWisdom")
}
