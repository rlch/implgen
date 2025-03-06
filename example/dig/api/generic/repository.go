package generic

type Entity[T any] struct {
	Field T
}

type Repository[T any] interface {
	Get(id string) (Entity[T], error)
	Gett(
		a string,
		b int, c bool,
	) (
		string,
		error,
	)
}

type NoGenericsRepository interface {
	Get(id string) (int, error)
}

type MultiGenericsRepository[A, B string, C float32] interface {
	A() (A, B, C)
}
