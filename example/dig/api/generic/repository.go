package generic

type Entity[T any] struct {
	Field T
}

type Repository[T any] interface {
	Get(id string) (Entity[T], error)
}

type NoGenericsRepository interface {
	Get(id string) (int, error)
}
