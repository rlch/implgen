package minecraft

import (
	"context"
	"time"
)

// Basic types
type Block struct {
	ID   string
	Type string
	Data map[string]any
}

type Player struct {
	Username string
	Health   float32
	XP       uint64
	Hunger   int
}

type Coordinates struct {
	X, Y, Z int
}

type (
	BlockType string
	ItemID    string
)

type CraftingRecipe struct {
	Name   string
	Inputs map[ItemID]int
	Output ItemID
}

// BlockRepository demonstrates type aliases and array types
type BlockRepository interface {
	// Type aliases
	PlaceBlock(ctx context.Context, blockType BlockType, pos Coordinates) error
	GetBlock(ctx context.Context, pos Coordinates) (*Block, error)

	// Array types and slices
	GetArea(ctx context.Context, from, to Coordinates) ([]Block, error)
	BulkPlace(ctx context.Context, blocks map[Coordinates]BlockType) error
}

// PlayerRepository demonstrates primitive type variations
type PlayerRepository interface {
	// Various numeric types
	UpdateHealth(ctx context.Context, username string, health float32) error
	AddXP(ctx context.Context, username string, xp uint64) error
	SetHunger(ctx context.Context, username string, hunger int) error

	// Pointers and nil returns
	FindPlayer(ctx context.Context, username string) (*Player, error)

	// Time types
	GetPlayTime(ctx context.Context, username string) (time.Duration, error)
}

// CraftingRepository demonstrates interface composition across packages
type CraftingRepository interface {
	// Map parameters and returns
	Craft(ctx context.Context, recipe CraftingRecipe, inventory map[ItemID]int) (map[ItemID]int, error)

	// Transfer operations
	TransferItems(ctx context.Context, fromPlayer, toPlayer string, items map[ItemID]int) error

	// Basic operation
	GetRecipes(ctx context.Context) ([]CraftingRecipe, error)
}
