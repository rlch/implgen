// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package genericimpl

import (
	"example/api/generic"

	"github.com/rotisserie/eris"
	"go.uber.org/dig"
)

type Dependencies[T any] struct {
	dig.In
	// Add dependencies here
}

func NewRepository[T any](deps Dependencies[T]) generic.Repository[T] {
	return &repositoryImpl[T]{
		Dependencies: deps,
	}
}

type repositoryImpl[T any] struct {
	Dependencies[T]
}

func (r *repositoryImpl[T]) Get(id string) (_ generic.Entity[T], err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "generic.Repository.Get")
		}
	}()
	panic("TODO: implement generic.Repository.Get")
}

func (r *repositoryImpl[T]) Gett(a string, b int, c bool) (_ string, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "generic.Repository.Gett")
		}
	}()
	panic("TODO: implement generic.Repository.Gett")
}
