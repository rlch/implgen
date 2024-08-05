// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate mockgen -source=../api/waltuh/repository.go -destination=waltuh/mocks/repository.go
//go:generate mockgen -source=../api/waltuh/another.go -destination=waltuh/mocks/another.go
//go:generate mockgen -source=../api/waltuh/repository.go -destination=waltuh/mocks/repository.go

import (
	waltuhimpl "example/internal/waltuh"

	"go.uber.org/fx"
)

var Repositories = fx.Options(
	waltuhimpl.Options,
	waltuhimpl.AnotherOptions,
	waltuhimpl.BOptions,
)
