package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

type IMetadata interface {
	Generate(imageIndex int, traits model.Traits) error
}

var _ IMetadata = &metadata{}

type metadata struct {
	collectionOutputDir string
}

func NewMetadata(collectionOutputDir string) IMetadata {
	return &metadata{
		collectionOutputDir: collectionOutputDir,
	}
}

func (s *metadata) Generate(imageIndex int, traits model.Traits) error {
	log.Debug().Msgf("Generating ERC721 metadata for image #%d", imageIndex)

	erc721Traits, err := traits.ToERC721()
	if err != nil {
		return errors.Annotate(err, "converting traits to ERC721 failed")
	}

	sortERCTraits(erc721Traits)

	b, err := json.MarshalIndent(erc721Traits, "", "	")
	if err != nil {
		return errors.Annotate(err, "marshalling json failed")
	}

	// Create traits file.
	fileName := filepath.Join(s.collectionOutputDir, fmt.Sprintf("%d.json", imageIndex))
	if err := os.WriteFile(fileName, b, 0777); err != nil {
		return errors.Annotate(err, "writing file failed")
	}

	return nil
}

// private

func sortERCTraits(ts []model.ERCTrait) {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].TraitType < ts[j].TraitType
	})
}
