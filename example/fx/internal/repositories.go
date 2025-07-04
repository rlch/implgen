// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=heisenberg/mocks.go -pkg=heisenbergimpl -rm -skip-ensure ../api/heisenberg ChemistryRepository MoneyRepository
//go:generate moq -out=spongebob/mocks.go -pkg=spongebobimpl -rm -skip-ensure ../api/spongebob JellyfishingRepository KrustyKrabRepository

import (
	"context"
	"example"
	"example/internal/heisenberg"
	"example/internal/spongebob"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/fx"

	_ "github.com/Southclaws/fault"

	_ "github.com/Southclaws/fault/fmsg"
)

var Repositories = fx.Options(
	heisenbergimpl.ChemistryOptions,
	heisenbergimpl.MoneyOptions,
	spongebobimpl.JellyfishingOptions,
	spongebobimpl.KrustyKrabOptions,
)
