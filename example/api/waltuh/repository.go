package waltuh

import (
	"context"
)

type Repository interface {
	MakeBreakfast(birthday, kilograms int) Waltuh
	SynthesizeMeth(ctx context.Context, flyPresent bool, withJesse bool) int
	MakeMoney(ctx context.Context, poundsOfMeth int) (int, error)
	DropWaltJrOffAtSchool(ctx context.Context) (bool, error)
	KillKrazy8(ctx context.Context, missingPlateShards int) (string, error)
	Get(ctx context.Context, id string) (string, error)
	Nope()
}

type BRepository interface{}
