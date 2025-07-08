// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package service

//go:generate moq -out=service/heisenberg/mocks.go -pkg=heisenbergimpl -rm -skip-ensure ../api/heisenberg NotificationService
//go:generate moq -out=service/spongebob/mocks.go -pkg=spongebobimpl -rm -skip-ensure ../api/spongebob FryService PattyService
import (
	heisenbergimpl "example/internal/service/heisenberg"
	spongebobimpl "example/internal/service/spongebob"

	"go.uber.org/fx"

	_ "github.com/Southclaws/fault"
	_ "github.com/Southclaws/fault/fmsg"
)

var Services = fx.Options(
	heisenbergimpl.NotificationOptions,
	spongebobimpl.FryOptions,
	spongebobimpl.PattyOptions,
)
