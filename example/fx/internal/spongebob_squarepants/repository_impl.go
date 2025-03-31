// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package spongebobsquarepantsimpl

import (
	spongebobsquarepants "example/api/spongebob_squarepants"

	"github.com/rotisserie/eris"
	"go.uber.org/fx"
)

type Dependencies struct {
	fx.In
	// Add dependencies here
}

var Options = fx.Options(
	fx.Provide(
		NewRepository,
	),
)

func NewRepository(deps Dependencies) spongebobsquarepants.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}

type repositoryImpl struct {
	Dependencies
}

func (r *repositoryImpl) Get(id string) (_ int, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebobsquarepants.Repository.Get")
		}
	}()
	panic("TODO: implement spongebobsquarepants.Repository.Get")
}
