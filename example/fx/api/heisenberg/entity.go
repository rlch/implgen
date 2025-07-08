package heisenberg

import "time"

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
