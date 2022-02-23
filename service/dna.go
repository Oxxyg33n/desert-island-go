package service

import (
	"github.com/Oxxyg33n/desert-island-go/repository"
	"github.com/juju/errors"
)

type IDNA interface {
	MarkAsExisting(dna string) error
}

var _ IDNA = &dna{}

type dna struct {
	dnaRepository repository.IDNA
}

func NewDNA(dnaRepository repository.IDNA) IDNA {
	return &dna{
		dnaRepository: dnaRepository,
	}
}

func (s *dna) MarkAsExisting(dna string) error {
	if s.dnaRepository.DNAExists(dna) {
		return errors.New("dna already exists")
	}

	s.dnaRepository.MarkAsExisting(dna)

	return nil
}
