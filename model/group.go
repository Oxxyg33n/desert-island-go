package model

import "github.com/google/uuid"

type Group struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Priority int       `json:"priority"`
	XPos     int       `json:"x_pos"`
	YPos     int       `json:"y_pos"`
}
