package spongebobsquarepants

type SpongebobSquarepantsRepository interface {
	Get(id string) (int, error)
}
