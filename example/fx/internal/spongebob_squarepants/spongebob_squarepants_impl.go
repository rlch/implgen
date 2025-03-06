// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package spongebobsquarepantsimpl

import (
	spongebobsquarepants "example/api/spongebob_squarepants"

	"github.com/rotisserie/eris"
	"go.uber.org/fx"
)

type SpongebobSquarepantsDependencies struct {
	fx.In
	// Add dependencies here
}

var SpongebobSquarepantsOptions = fx.Options(
	fx.Provide(
		NewSpongebobSquarepantsRepository,
	),
)

func NewSpongebobSquarepantsRepository(deps SpongebobSquarepantsDependencies) spongebobsquarepants.SpongebobSquarepantsRepository {
	return &spongebobSquarepantsRepositoryImpl{
		SpongebobSquarepantsDependencies: deps,
	}
}

type spongebobSquarepantsRepositoryImpl struct {
	SpongebobSquarepantsDependencies
}

func (r *spongebobSquarepantsRepositoryImpl) Get(id string) (_ int, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebobsquarepants.SpongebobSquarepantsRepository.Get")
		}
	}()
	panic("TODO: implement spongebobsquarepants.SpongebobSquarepantsRepository.Get")
}
