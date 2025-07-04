// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package minecraftimpl

import (
	"context"
	"time"

	"example/api/minecraft"

	"github.com/rotisserie/eris"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/dig"
)

type PlayerDependencies struct {
	dig.In
	// Add dependencies here
}

func NewPlayerRepository(deps PlayerDependencies) minecraft.PlayerRepository {
	return &playerRepositoryImpl{
		PlayerDependencies: deps,
	}
}

type playerRepositoryImpl struct {
	PlayerDependencies
}

func (r *playerRepositoryImpl) UpdateHealth(ctx context.Context, username string, health float32) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Player.UpdateHealth")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.PlayerRepository.UpdateHealth")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.PlayerRepository.UpdateHealth")
}

func (r *playerRepositoryImpl) AddXP(ctx context.Context, username string, xp uint64) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Player.AddXP")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.PlayerRepository.AddXP")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.PlayerRepository.AddXP")
}

func (r *playerRepositoryImpl) SetHunger(ctx context.Context, username string, hunger int) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Player.SetHunger")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.PlayerRepository.SetHunger")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.PlayerRepository.SetHunger")
}

func (r *playerRepositoryImpl) FindPlayer(ctx context.Context, username string) (_ *minecraft.Player, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Player.FindPlayer")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.PlayerRepository.FindPlayer")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.PlayerRepository.FindPlayer")
}

func (r *playerRepositoryImpl) GetPlayTime(ctx context.Context, username string) (_ time.Duration, err error) {
	ctx, span := otel.GetTracerProvider().Tracer("minecraft").Start(ctx, "Player.GetPlayTime")
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "minecraft.PlayerRepository.GetPlayTime")
			span.SetStatus(codes.Error, "")
			span.RecordError(err)
		}
		span.End()
	}()
	_ = ctx
	panic("TODO: implement minecraft.PlayerRepository.GetPlayTime")
}
