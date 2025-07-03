package shrek

import (
	"context"
	"io"
)

// SwampRepository demonstrates basic operations
type SwampRepository interface {
	// Basic CRUD
	CleanSwamp(ctx context.Context) error
	GetSwampStatus(ctx context.Context) (string, error)

	// Character operations
	EvictIntruders(ctx context.Context, characters []string) error
	WelcomeVisitor(ctx context.Context, visitor string) error
}

// QuestRepository demonstrates interface{} usage
type QuestRepository interface {
	// Any type parameters
	StartQuest(ctx context.Context, hero any, objective any) (string, error)

	// Basic operations
	RegisterCharacter(ctx context.Context, name string, charType string) error

	// Map returns
	GetQuestStatus(ctx context.Context, questID string) (map[string]any, error)
}

// AdvancedRepository demonstrates complex signatures
type AdvancedRepository interface {
	// Variadic parameters
	FormParty(ctx context.Context, leader string, members ...string) error

	// Function type parameters
	ProcessCharacters(ctx context.Context, filter func(string) bool) ([]string, error)

	// io interfaces
	ExportStory(ctx context.Context, questID string, writer io.Writer) error
	ImportLegend(ctx context.Context, reader io.Reader) error

	// No parameters, no context
	GetMagicMirrorWisdom() string
}
