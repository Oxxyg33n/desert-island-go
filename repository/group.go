package repository

import (
	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/google/uuid"
	"github.com/juju/errors"
)

type IGroup interface {
	Create(trait model.Group) model.Group
	GetByID(id uuid.UUID) (*model.Group, error)
	GetByName(name string) (*model.Group, error)
	GetAll() ([]model.Group, error)
}

var _ IGroup = &group{}

type group struct {
	groupByIDMap   map[uuid.UUID]model.Group
	groupByNameMap map[string]model.Group
}

func NewGroup() *group {
	return &group{
		groupByIDMap:   make(map[uuid.UUID]model.Group),
		groupByNameMap: make(map[string]model.Group),
	}
}

func (r *group) Create(group model.Group) model.Group {
	r.groupByIDMap[group.ID] = group
	r.groupByNameMap[group.Name] = group

	return group
}

func (r *group) GetByID(groupID uuid.UUID) (*model.Group, error) {
	foundGroup, ok := r.groupByIDMap[groupID]
	if !ok {
		return nil, errors.New("getting group by ID failed")
	}

	return &foundGroup, nil
}

func (r *group) GetByName(name string) (*model.Group, error) {
	foundGroup, ok := r.groupByNameMap[name]
	if !ok {
		return nil, errors.New("getting group by name failed")
	}

	return &foundGroup, nil
}

func (r *group) GetAll() ([]model.Group, error) {
	if len(r.groupByIDMap) == 0 {
		return nil, errors.New("no groups found")
	}

	var groups []model.Group
	for _, g := range r.groupByIDMap {
		groups = append(groups, g)
	}

	return groups, nil
}
