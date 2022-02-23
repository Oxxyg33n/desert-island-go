package service

import (
	"fmt"

	"github.com/Oxxyg33n/desert-island-go/repository"
)

type IDNA interface {
	DNAExists(dna string) bool
	MarkAsExisting(dna string) error
}

var _ IDNA = &dna{}

type dna struct {
	dnaRepository repository.IDNA
	dnaPrefix     string
}

func NewDNA(dnaRepository repository.IDNA, dnaPrefix string) IDNA {
	return &dna{
		dnaRepository: dnaRepository,
		dnaPrefix:     dnaPrefix,
	}
}

func (s *dna) DNAExists(dna string) bool {
	dna = fmt.Sprintf("%s-%s", s.dnaPrefix, dna)

	return s.dnaRepository.DNAExists(dna)
}

func (s *dna) MarkAsExisting(dna string) error {
	dna = fmt.Sprintf("%s-%s", s.dnaPrefix, dna)

	s.dnaRepository.MarkAsExisting(dna)

	return nil
}
