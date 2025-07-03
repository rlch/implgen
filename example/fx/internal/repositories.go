// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=heisenberg/mocks.go -pkg=heisenbergimpl -rm -skip-ensure ../api/heisenberg ChemistryRepository MoneyRepository
//go:generate moq -out=spongebob/mocks.go -pkg=spongebobimpl -rm -skip-ensure ../api/spongebob JellyfishingRepository KrustyKrabRepository

import (
	heisenbergimpl "example/internal/heisenberg"
	spongebobimpl "example/internal/spongebob"
)

var RepositoryFactories = []any{
	heisenbergimpl.NewChemistryRepository,
	heisenbergimpl.NewMoneyRepository,
	spongebobimpl.NewJellyfishingRepository,
	spongebobimpl.NewKrustyKrabRepository,
}
