package heisenbergimpl

import (
	"example/api/heisenberg"

	"github.com/rotisserie/eris"
)

func (r *chemistryRepositoryImpl) OptimizeFormula(formula heisenberg.Formula) (_ heisenberg.Formula, _ []string, err error) {
	defer func() {
		if err != nil {
			err = eris.Wrap(err, "heisenberg.ChemistryRepository.OptimizeFormula")
		}
	}()
	panic("TODO: implement heisenberg.ChemistryRepository.OptimizeFormula")
}
