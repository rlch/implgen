// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=waltuh/mocks.go -pkg=waltuhimpl -rm -skip-ensure ../api/waltuh AnotherRepository BRepository Repository

import (
	"go.uber.org/fx"

	waltuhimpl "example/internal/waltuh"
)

var Repositories = fx.Options(
	waltuhimpl.Options,
	waltuhimpl.AnotherOptions,
	waltuhimpl.BOptions,
)
