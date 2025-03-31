// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=waltuh/nested/mocks.go -pkg=nestedimpl -rm -skip-ensure ../api/waltuh/nested Repository
//go:generate moq -out=spongebob_squarepants/mocks.go -pkg=spongebobsquarepantsimpl -rm -skip-ensure ../api/spongebob_squarepants Repository
//go:generate moq -out=waltuh/mocks.go -pkg=waltuhimpl -rm -skip-ensure ../api/waltuh AnotherRepository BRepository Repository

import (
	spongebobsquarepantsimpl "example/internal/spongebob_squarepants"
	waltuhimpl "example/internal/waltuh"
	nestedimpl "example/internal/waltuh/nested"

	"go.uber.org/fx"
)

var Repositories = fx.Options(
	nestedimpl.Options,
	spongebobsquarepantsimpl.Options,
	waltuhimpl.Options,
	waltuhimpl.AnotherOptions,
	waltuhimpl.BOptions,
)
