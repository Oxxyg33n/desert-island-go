package model

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"

	"github.com/google/uuid"
	"github.com/juju/errors"
)

type Traits []Trait

type Trait struct {
	ID       uuid.UUID
	Group    TraitGroup
	Name     string
	Image    []byte
	Rareness Rareness
	DNAIndex int
}

func (t *Trait) ToImageLayer() (*ImageLayer, error) {
	if t == nil {
		return nil, errors.New("trait is nil")
	}

	img, err := png.Decode(bytes.NewReader(t.Image))
	if err != nil {
		return nil, errors.Annotate(err, "decoding png failed")
	}

	return &ImageLayer{
		Image:    img,
		Priority: t.Group.Priority,
	}, nil
}

func (ts Traits) ToImageLayers() ([]ImageLayer, error) {
	if ts == nil {
		return nil, errors.New("traits is nil")
	}

	layers := make([]ImageLayer, len(ts))
	for i, trait := range ts {
		layer, err := trait.ToImageLayer()
		if err != nil {
			return nil, errors.Annotate(err, "converting image to layer failed")
		}

		layers[i] = *layer
	}

	return layers, nil
}

func (ts Traits) ToDNA() string {
	var dna string
	for _, trait := range ts {
		dna += fmt.Sprintf("%d-", trait.DNAIndex)
	}

	dna = strings.TrimSuffix(dna, "-")

	return dna
}

func (t *Trait) ToERC721() (*ERCTrait, error) {
	if t == nil {
		return nil, errors.New("trait is nil")
	}

	return &ERCTrait{
		TraitType: strings.ToUpper(t.Group.Name),
		Value:     strings.ToUpper(t.Name),
	}, nil
}

func (ts Traits) ToERC721() ([]ERCTrait, error) {
	if ts == nil {
		return nil, errors.New("traits is nil")
	}

	erc721Traits := make([]ERCTrait, len(ts))
	for i, trait := range ts {
		erc721Trait, err := trait.ToERC721()
		if err != nil {
			return nil, errors.Annotate(err, "converting trait to ERC721 metadata failed")
		}

		erc721Traits[i] = *erc721Trait
	}

	return erc721Traits, nil
}
