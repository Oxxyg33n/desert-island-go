package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Oxxyg33n/desert-island-go/configuration"
	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

type IMetadata interface {
	Generate(imageIndex int, traits model.Traits) error
}

var _ IMetadata = &metadata{}

type metadata struct {
	cfg configuration.Configuration
}

func NewMetadata(cfg configuration.Configuration) IMetadata {
	return &metadata{
		cfg: cfg,
	}
}

func (s *metadata) Generate(imageIndex int, traits model.Traits) error {
	log.Debug().Msgf("Generating ERC721 metadata for image #%d", imageIndex)

	erc721Traits, err := traits.ToERC721()
	if err != nil {
		return errors.Annotate(err, "converting traits to ERC721 failed")
	}

	sortERCTraits(erc721Traits)

	ercMetadata := model.NewERCMetadata(
		s.cfg.CollectionName, s.cfg.CollectionDescription, "",
		imageIndex,
		erc721Traits,
	)

	b, err := json.MarshalIndent(ercMetadata, "", "	")
	if err != nil {
		return errors.Annotate(err, "marshalling json failed")
	}

	// Create traits file.
	fileName := filepath.Join(s.cfg.CollectionOutputDir, "traits", fmt.Sprintf("%d.json", imageIndex))
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
