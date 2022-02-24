package model

import "github.com/google/uuid"

type TraitGroup struct {
	ID         uuid.UUID
	Name       string
	CanSkip    bool
	SkipChance float32
	Priority   int
}
