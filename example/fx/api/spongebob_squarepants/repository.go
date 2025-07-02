package spongebobsquarepants

type Repository interface {
	Get(id string) (int, error)
	Get2(id string) (int, error)
}
