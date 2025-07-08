package spongebob

import (
	"context"
)

// PattyService handles Krabby Patty operations
type PattyService interface {
	// Grill a new patty
	GrillPatty(ctx context.Context, orderID string) error
	
	// Check if patty is ready
	IsPattyReady(ctx context.Context, orderID string) (bool, error)
	
	// Serve the patty
	ServePatty(ctx context.Context, orderID string, customerName string) error
}

// FryService manages the fry station
type FryService interface {
	// Start cooking fries
	CookFries(ctx context.Context, size string) (string, error)
	
	// Season the fries
	SeasonFries(ctx context.Context, batchID string, seasoning string) error
}