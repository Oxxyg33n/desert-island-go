package model

import "github.com/google/uuid"

type Group struct {
	ID       uuid.UUID
	Name     string
	Priority int
	XPos     int
	YPos     int
}
