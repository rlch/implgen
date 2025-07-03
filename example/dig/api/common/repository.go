package common

import "context"

// Healthable defines basic health operations
type Healthable interface {
	IsHealthy(ctx context.Context) error
}

// BaseRepository embeds Healthable and adds metrics
type BaseRepository interface {
	Healthable
	GetMetrics(ctx context.Context) (map[string]int64, error)
}