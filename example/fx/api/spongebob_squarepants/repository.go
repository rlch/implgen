package spongebobsquarepants

type Repository interface {
	Get(id string) (int, error)
}

type LessonDocumentationReconciler[N any, PN interface {
	Repository
	*N
}] interface {
	Reconcile() error
}
