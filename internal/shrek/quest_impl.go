// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package shrekimpl

import (
	"context"

	"example/api/shrek"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/dig"
)

type QuestDependencies struct {
	dig.In
	// Add dependencies here
}

func NewQuestRepository(deps QuestDependencies) shrek.QuestRepository {
	return &questRepositoryImpl{
		QuestDependencies: deps,
	}
}

type questRepositoryImpl struct {
	QuestDependencies
}

func (r *questRepositoryImpl) StartQuest(ctx context.Context, hero, objective any) (_ string, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Quest.StartQuest")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.QuestRepository.StartQuest")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.QuestRepository.StartQuest")
}

func (r *questRepositoryImpl) RegisterCharacter(ctx context.Context, name, charType string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Quest.RegisterCharacter")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.QuestRepository.RegisterCharacter")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.QuestRepository.RegisterCharacter")
}

func (r *questRepositoryImpl) GetQuestStatus(ctx context.Context, questID string) (_ map[string]any, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("shrek").Start(ctx, "Quest.GetQuestStatus")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "shrek.QuestRepository.GetQuestStatus")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement shrek.QuestRepository.GetQuestStatus")
}
