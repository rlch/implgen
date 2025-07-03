package spongebob

import "context"

// Basic types  
type KrabbyPatty struct {
	ID       string
	Toppings []string
	Quality  int
}

type MenuItem struct {
	Name  string
	Price float64
}

type Jellyfish interface {
	Sting() int
	GetColor() string
}

// KrustyKrabRepository demonstrates basic repository patterns
type KrustyKrabRepository interface {
	// Basic CRUD
	CookPatty(ctx context.Context, toppings []string) (*KrabbyPatty, error)
	ServeCustomer(ctx context.Context, patty KrabbyPatty, customerName string) error
	
	// Menu operations
	GetMenu(ctx context.Context) ([]MenuItem, error)
	UpdateMenu(ctx context.Context, items []MenuItem) error
}

// JellyfishingRepository shows interface types and slices
type JellyfishingRepository interface {
	// Interface parameters
	CatchJellyfish(ctx context.Context, jellyfish Jellyfish) error
	ReleaseJellyfish(ctx context.Context, jellyfish []Jellyfish) error
	
	// Basic operations
	CountJellyfish(ctx context.Context, area string) (int, error)
	GetBestSpot(ctx context.Context) (string, error)
}