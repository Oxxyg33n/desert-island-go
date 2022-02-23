package repository

type IDNA interface {
	DNAExists(dna string) bool
	MarkAsExisting(dna string)
}

var _ IDNA = &dnaRepository{}

type dnaRepository struct {
	dnaMap map[string]struct{}
}

func NewDNA() *dnaRepository {
	return &dnaRepository{
		dnaMap: make(map[string]struct{}),
	}
}

func (r *dnaRepository) DNAExists(dna string) bool {
	if _, ok := r.dnaMap[dna]; ok {
		return true
	}

	return false
}

func (r *dnaRepository) MarkAsExisting(dna string) {
	r.dnaMap[dna] = exists
}

// private

var exists = struct{}{}
