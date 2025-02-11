// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package nestedimpl

import (
	"example/api/waltuh/nested"

	"go.uber.org/dig"
)

type Dependencies struct {
	dig.In
	// Add dependencies here
}

func NewRepository(deps Dependencies) nested.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}

type repositoryImpl struct {
	Dependencies
}

func (r *repositoryImpl) HelloWorld() string {
	panic("TODO: implement nested.Repository.HelloWorld")
}
