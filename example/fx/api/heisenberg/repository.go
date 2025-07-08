package heisenberg

import (
	"context"
	"time"
)

// Simple types for testing
type Formula struct {
	Name   string
	Purity float64
}

type Batch struct {
	ID       string
	Formula  Formula
	Quantity float64
	CookedAt time.Time
}

// ChemistryRepository demonstrates basic patterns
type ChemistryRepository interface {
	// Basic CRUD
	Cook(ctx context.Context, formula Formula) (*Batch, error)
	GetBatch(ctx context.Context, id string) (*Batch, error)

	// Multiple return values
	OptimizeFormula(formula Formula) (Formula, []string, error)
	
	// New method to test generation
	TestMethod(ctx context.Context, input string) (string, error)
}

// MoneyRepository demonstrates different types
type MoneyRepository interface {
	// Nullable types
	Launder(ctx context.Context, amount *float64) (*float64, error)

	// Slices
	ProcessPayments(ctx context.Context, amounts []float64) ([]string, error)
}

