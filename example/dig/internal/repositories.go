// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=minecraft/mocks.go -pkg=minecraftimpl -rm -skip-ensure ../api/minecraft BlockRepository CraftingRepository PlayerRepository
//go:generate moq -out=shrek/mocks.go -pkg=shrekimpl -rm -skip-ensure ../api/shrek AdvancedRepository QuestRepository SwampRepository

import (
	minecraftimpl "example/internal/minecraft"
	shrekimpl "example/internal/shrek"
)

var RepositoryFactories = []any{
	minecraftimpl.NewBlockRepository,
	minecraftimpl.NewCraftingRepository,
	minecraftimpl.NewPlayerRepository,
	shrekimpl.NewAdvancedRepository,
	shrekimpl.NewQuestRepository,
	shrekimpl.NewSwampRepository,
}
