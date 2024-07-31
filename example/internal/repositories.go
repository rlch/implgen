// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate mockgen -source=api/waltuh/repository.go -destination=internal/waltuh/mocks/repository.go

import (
	waltuhimpl "example/internal/waltuh"

	"go.uber.org/fx"
)

var Repositories = fx.Options(
	waltuhimpl.Options,
	waltuhimpl.BOptions,
)
