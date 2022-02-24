package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/Oxxyg33n/desert-island-go/repository"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

func init() {
	// Seed randomizer on each run
	rand.Seed(time.Now().UnixNano())
}

type ITrait interface {
	Import(root string) error
	GetRandomTraits() (model.Traits, error)
}

var _ ITrait = &trait{}

type trait struct {
	groupRepository repository.IGroup
	traitRepository repository.ITrait

	skipMultipleLayers bool

	groupCfgByName map[string]model.GroupConfiguration

	// Width and height that will be used to resize images
	imageWidth, imageHeight int

	// Traits grouped by group
	groupedTraits map[model.TraitGroup][]model.Trait
}

// NewTrait creates new Trait service and loads priority configuration from JSON file.
func NewTrait(
	traitRepository repository.ITrait,
	groupRepository repository.IGroup,
	imageWidth, imageHeight int,
) ITrait {
	s := &trait{
		traitRepository: traitRepository,
		groupRepository: groupRepository,
		imageWidth:      imageWidth,
		imageHeight:     imageHeight,

		groupedTraits:  make(map[model.TraitGroup][]model.Trait),
		groupCfgByName: make(map[string]model.GroupConfiguration),
	}

	s.loadCollectionConfiguration()

	return s
}

func (s *trait) Import(root string) error {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return errors.Annotate(err, "finding directory failed")
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if info.IsDir() {
			return nil
		}

		fileExtension := filepath.Ext(info.Name())

		if fileExtension != ".png" {
			return errors.New("file is not a .png")
		}

		splitted := strings.Split(path, "/")
		if len(splitted) < 3 {
			return errors.New("invalid directories format, expected at least 3 elements")
		}

		groupName := splitted[1]
		traitRarity := splitted[2]
		traitName := strings.TrimSuffix(info.Name(), fileExtension)

		group, err := s.groupRepository.GetByName(groupName)
		if err != nil {
			log.Debug().Msgf("Trait group not found, creating new one: group=%s", groupName)
			newGroup := s.groupRepository.Create(model.TraitGroup{
				ID:         uuid.New(),
				CanSkip:    s.groupCfgByName[groupName].CanSkip,
				SkipChance: s.groupCfgByName[groupName].SkipChance,
				Name:       groupName,
				Priority:   s.groupCfgByName[groupName].Priority,
			})

			group = &newGroup
		}

		buf := bytes.Buffer{}
		src, err := imaging.Open(path)
		if err != nil {
			return errors.Annotate(err, "opening image failed")
		}

		src = imaging.Resize(src, s.imageWidth, s.imageHeight, imaging.Lanczos)

		if err := imaging.Encode(&buf, src, imaging.PNG); err != nil {
			return errors.Annotate(err, "encoding image to png failed")
		}

		if err := s.traitRepository.Create(model.Trait{
			Name:     traitName,
			Group:    *group,
			Image:    buf.Bytes(),
			Rareness: model.Rareness(traitRarity),
		}); err != nil {
			return errors.Annotate(err, "creating trait failed")
		}

		return nil
	}); err != nil {
		return errors.Annotate(err, "walking path failed")
	}

	groups, err := s.groupRepository.GetAll()
	if err != nil {
		return errors.Annotate(err, "getting all groups failed")
	}

	var strGroups string
	for _, group := range groups {
		strGroups += fmt.Sprintf("%s ", group.Name)
	}

	log.Debug().Msgf("%d layer groups created: %s", len(groups), strGroups)

	for _, group := range groups {
		traits, err := s.traitRepository.GetByGroupID(group.ID)
		if err != nil {
			return errors.Annotate(err, "getting traits by group ID failed")
		}

		if len(traits) == 0 {
			return errors.Errorf("no traits found for group: group=%s", group.Name)
		}

		s.groupedTraits[group] = traits
	}

	return nil
}

func (s *trait) GetRandomTraits() (model.Traits, error) {
	mustSkip := false
	randomTraits := make([]model.Trait, 0)
	for group, traits := range s.groupedTraits {
		if !mustSkip && s.skipGroup(group) {
			if !s.skipMultipleLayers {
				mustSkip = true
			}

			continue
		}

		t, err := getRandomTrait(traits)
		if err != nil {
			return nil, errors.Annotate(err, "getting random trait failed")
		}

		randomTraits = append(randomTraits, *t)
	}

	return randomTraits, nil
}

// private

func (s *trait) skipGroup(group model.TraitGroup) bool {
	// Check if group exists in configuration,
	// if not then skip the group.
	groupConfiguration, ok := s.groupCfgByName[group.Name]
	if !ok {
		return true
	}

	// Check if group can be skipped,
	// if not then do not skip layer
	if canSkip := groupConfiguration.CanSkip; !canSkip {
		return false
	}

	chance := groupConfiguration.SkipChance
	// Generate random float32 in min 0, max 100 range and check it against skipping chance
	r := rand.Float32() * 100 // 100% maximum

	return r <= chance
}

func getRandomTrait(traits []model.Trait) (*model.Trait, error) {
	pdf, err := getProbabilityDensityVector(traits)
	if err != nil {
		return nil, errors.Annotate(err, "getting probability density vector failed")
	}

	// get cdf
	traitsLen := len(traits)
	cdf := make([]float32, traitsLen)
	cdf[0] = pdf[0]
	for i := 1; i < traitsLen; i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}

	random := sample(cdf)
	if !(traitsLen > random) {
		return nil, errors.Errorf("generated trait out of range: max size=%d, generated index=%d", traitsLen, random)
	}

	return &traits[random], nil
}

// private

func getProbabilityDensityVector(traits []model.Trait) ([]float32, error) {
	traitsLen := len(traits)
	probabilityVector := make([]float32, traitsLen)

	var (
		commonChance = float32(100/traitsLen) / 100
		rareChance   = commonChance / 3
		epicChance   = commonChance / 4
	)

	chanceOffset := float32(1.00)
	commonCounter := float32(0)

	for i, t := range traits {
		switch t.Rareness {
		case model.RarenessRare:
			probabilityVector[i] = rareChance
			chanceOffset -= rareChance

		case model.RarenessEpic:
			probabilityVector[i] = epicChance
			chanceOffset -= epicChance

		default:
			commonCounter++
		}
	}

	for i, p := range probabilityVector {
		if p == 0 {
			probabilityVector[i] = chanceOffset / commonCounter
		}
	}

	if err := checkProbabilityVector(probabilityVector); err != nil {
		return nil, errors.Annotate(err, "checking probability vector failed")
	}

	return probabilityVector, nil
}

func checkProbabilityVector(vector []float32) error {
	sum := float32(0)
	checkSum := float32(1)

	for _, p := range vector {
		sum += p
	}

	if !(sum >= checkSum-0.1 || sum >= checkSum+0.1) {
		return errors.Errorf("expected probability vector checksum %.2f but got %.2f", checkSum, sum)
	}

	return nil
}

func sample(cdf []float32) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Float32()
	bucket := 0

	for r > cdf[bucket] {
		bucket++
	}

	return bucket
}

func (s *trait) loadCollectionConfiguration() {
	b, err := os.ReadFile("collection_configuration.json")
	if err != nil {
		log.Fatal().Msg(errors.Annotate(err, "reading collection configuration file failed").Error())
	}

	var collectionConfiguration model.CollectionConfiguration
	if err := json.Unmarshal(b, &collectionConfiguration); err != nil {
		log.Fatal().Msg(errors.Annotate(err, "unmarshalling json failed").Error())
	}

	s.groupCfgByName = make(map[string]model.GroupConfiguration, len(collectionConfiguration.Layers.Groups))
	for _, cfg := range collectionConfiguration.Layers.Groups {
		s.groupCfgByName[cfg.Name] = cfg
	}

	s.skipMultipleLayers = collectionConfiguration.Layers.SkipMultiple
}
