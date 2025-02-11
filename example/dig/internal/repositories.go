// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=generic/mocks.go -pkg=genericimpl -rm -skip-ensure ../api/generic MultiGenericsRepository NoGenericsRepository Repository
//go:generate moq -out=waltuh/nested/mocks.go -pkg=nestedimpl -rm -skip-ensure ../api/waltuh/nested Repository
//go:generate moq -out=waltuh/mocks.go -pkg=waltuhimpl -rm -skip-ensure ../api/waltuh AnotherRepository BRepository Repository

import (
	genericimpl "example/internal/generic"
	nestedimpl "example/internal/waltuh/nested"

	waltuhimpl "example/internal/waltuh"
)

var RepositoryFactories = []any{

	genericimpl.NewNoGenericsRepository,
	nestedimpl.NewRepository,
	waltuhimpl.NewRepository,
	waltuhimpl.NewAnotherRepository,
	waltuhimpl.NewBRepository,
}
