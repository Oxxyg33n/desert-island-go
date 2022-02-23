package repository

import (
	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/google/uuid"
	"github.com/juju/errors"
)

type ITrait interface {
	Create(trait model.Trait) error
	GetByID(id uuid.UUID) (*model.Trait, error)
	GetByGroupID(id uuid.UUID) ([]model.Trait, error)
}

var _ ITrait = &Trait{}

type Trait struct {
	traitByIDMap          map[uuid.UUID]model.Trait
	traitByNameMap        map[string]model.Trait
	traitsByGroupIDMap    map[uuid.UUID][]model.Trait
	lastTraitIndexByGroup map[string]int
}

func NewTrait() *Trait {
	return &Trait{
		traitByIDMap:          make(map[uuid.UUID]model.Trait),
		traitByNameMap:        make(map[string]model.Trait),
		traitsByGroupIDMap:    make(map[uuid.UUID][]model.Trait),
		lastTraitIndexByGroup: make(map[string]int),
	}
}

func (r *Trait) Create(trait model.Trait) error {
	r.traitByIDMap[trait.ID] = trait
	r.traitByNameMap[trait.Name] = trait
	r.lastTraitIndexByGroup[trait.Group.Name]++

	if _, ok := r.traitsByGroupIDMap[trait.Group.ID]; !ok {
		r.traitsByGroupIDMap[trait.Group.ID] = make([]model.Trait, 0)
	}

	r.traitsByGroupIDMap[trait.Group.ID] = append(r.traitsByGroupIDMap[trait.Group.ID], trait)

	return nil
}

func (r *Trait) GetByID(traitID uuid.UUID) (*model.Trait, error) {
	t, ok := r.traitByIDMap[traitID]
	if !ok {
		return nil, errors.New("getting Trait by ID failed")
	}

	return &t, nil
}

func (r *Trait) GetByGroupID(groupID uuid.UUID) ([]model.Trait, error) {
	return r.traitsByGroupIDMap[groupID], nil
}
