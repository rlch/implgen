package waltuh

import "context"

type AnotherRepository interface {
	A() (string, error)
	B(ctx context.Context) error
	C() (string, error)
}
