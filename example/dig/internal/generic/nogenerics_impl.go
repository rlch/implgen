// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package genericimpl

import (
	"example/api/generic"

	"github.com/rotisserie/eris"
	"go.uber.org/dig"
)

type NoGenericsDependencies struct {
	dig.In
	// Add dependencies here
}

func NewNoGenericsRepository(deps NoGenericsDependencies) generic.NoGenericsRepository {
	return &noGenericsRepositoryImpl{
		NoGenericsDependencies: deps,
	}
}

type noGenericsRepositoryImpl struct {
	NoGenericsDependencies
}

func (r *noGenericsRepositoryImpl) Get(id string) (_ int, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "generic.NoGenericsRepository.Get")
		}
	}()
	panic("TODO: implement generic.NoGenericsRepository.Get")
}
