// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package genericimpl

import (
	"example/api/generic"

	"go.uber.org/dig"
)

type MultiGenericsDependencies[A, B string, C float32] struct {
	dig.In
	// Add dependencies here
}

func NewMultiGenericsRepository[A, B string, C float32](deps MultiGenericsDependencies[A, B, C]) generic.MultiGenericsRepository[A, B, C] {
	return &multiGenericsRepositoryImpl[A, B, C]{
		MultiGenericsDependencies: deps,
	}
}

type multiGenericsRepositoryImpl[A, B string, C float32] struct {
	MultiGenericsDependencies[A, B, C]
}

func (r *multiGenericsRepositoryImpl[A, B, C]) A() (A, B, C) {
	panic("TODO: implement generic.MultiGenericsRepository.A")
}
