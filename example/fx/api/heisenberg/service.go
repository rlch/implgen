package heisenberg

import (
	"context"
)

// NotificationService demonstrates service pattern
type NotificationService interface {
	// Send notification
	SendAlert(ctx context.Context, message string) error
	
	// Get notification status
	GetStatus(ctx context.Context, id string) (string, error)
}