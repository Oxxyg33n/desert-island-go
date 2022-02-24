package repository

import (
	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/google/uuid"
	"github.com/juju/errors"
)

type IGroup interface {
	Create(trait model.TraitGroup) model.TraitGroup
	GetByID(id uuid.UUID) (*model.TraitGroup, error)
	GetByName(name string) (*model.TraitGroup, error)
	GetAll() ([]model.TraitGroup, error)
}

var _ IGroup = &group{}

type group struct {
	groupByIDMap   map[uuid.UUID]model.TraitGroup
	groupByNameMap map[string]model.TraitGroup
}

func NewGroup() IGroup {
	return &group{
		groupByIDMap:   make(map[uuid.UUID]model.TraitGroup),
		groupByNameMap: make(map[string]model.TraitGroup),
	}
}

func (r *group) Create(group model.TraitGroup) model.TraitGroup {
	r.groupByIDMap[group.ID] = group
	r.groupByNameMap[group.Name] = group

	return group
}

func (r *group) GetByID(groupID uuid.UUID) (*model.TraitGroup, error) {
	foundGroup, ok := r.groupByIDMap[groupID]
	if !ok {
		return nil, errors.New("getting group by ID failed")
	}

	return &foundGroup, nil
}

func (r *group) GetByName(name string) (*model.TraitGroup, error) {
	foundGroup, ok := r.groupByNameMap[name]
	if !ok {
		return nil, errors.New("getting group by name failed")
	}

	return &foundGroup, nil
}

func (r *group) GetAll() ([]model.TraitGroup, error) {
	if len(r.groupByIDMap) == 0 {
		return nil, errors.New("no groups found")
	}

	var groups []model.TraitGroup
	for _, g := range r.groupByIDMap {
		groups = append(groups, g)
	}

	return groups, nil
}
