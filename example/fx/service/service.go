// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package service

//go:generate moq -out=heisenberg/mocks.go -pkg=heisenbergimpl -rm -skip-ensure ../api/heisenberg NotificationService
//go:generate moq -out=spongebob/mocks.go -pkg=spongebobimpl -rm -skip-ensure ../api/spongebob FryService PattyService
import (
	heisenbergimpl "example/service/heisenberg"
	spongebobimpl "example/service/spongebob"

	"go.uber.org/fx"

	_ "github.com/Southclaws/fault"
	_ "github.com/Southclaws/fault/fmsg"
)

var Services = fx.Options(
	heisenbergimpl.NotificationOptions,
	spongebobimpl.FryOptions,
	spongebobimpl.PattyOptions,
)
