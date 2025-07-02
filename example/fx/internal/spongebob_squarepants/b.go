package spongebobsquarepantsimpl

import "github.com/rotisserie/eris"

func (r *repositoryImpl) Get2(id string) (_ int, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "spongebobsquarepants.Repository.Get2")
		}
	}()
	panic("TODO: implement spongebobsquarepants.Repository.Get2")
}
