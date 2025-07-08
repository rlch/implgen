// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package repository

//go:generate moq -out=repository/heisenberg/mocks.go -pkg=heisenbergimpl -rm -skip-ensure ../api/heisenberg ChemistryRepository MoneyRepository
//go:generate moq -out=repository/spongebob/mocks.go -pkg=spongebobimpl -rm -skip-ensure ../api/spongebob JellyfishingRepository KrustyKrabRepository
import (
	heisenbergimpl "example/internal/repository/heisenberg"
	spongebobimpl "example/internal/repository/spongebob"

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
