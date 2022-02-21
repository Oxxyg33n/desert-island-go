package model

import (
	"bytes"
	"image/png"

	"github.com/google/uuid"
	"github.com/juju/errors"
)

type Trait struct {
	ID       uuid.UUID `json:"id"`
	Group    Group     `json:"group"`
	Name     string    `json:"name"`
	Image    []byte    `json:"image"`
	Rareness Rareness  `json:"rareness"`
}

func (t *Trait) ToImageLayer() (*ImageLayer, error) {
	if t == nil {
		return nil, nil
	}

	img, err := png.Decode(bytes.NewReader(t.Image))
	if err != nil {
		return nil, errors.Annotate(err, "decoding png failed")
	}

	return &ImageLayer{
		Image:    img,
		Priority: t.Group.Priority,
		XPos:     t.Group.XPos,
		YPos:     t.Group.YPos,
	}, nil
}
