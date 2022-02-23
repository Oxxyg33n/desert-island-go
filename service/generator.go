package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Oxxyg33n/desert-island-go/configuration"
	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

type IGenerator interface {
	Generate(imageIndex int) error
}

var _ IGenerator = &generator{}

type generator struct {
	cfg          configuration.Configuration
	traitService ITrait
	dnaService   IDNA
}

func NewGenerator(cfg configuration.Configuration, traitService ITrait, dnaService IDNA) *generator {
	if _, err := os.Stat(cfg.CollectionOutputDir); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.CollectionOutputDir, 0777); err != nil {
			log.Fatal().Msg("creating output directory failed")
		}
	}

	return &generator{
		cfg:          cfg,
		traitService: traitService,
		dnaService:   dnaService,
	}
}

func (s *generator) Generate(imageIndex int) error {
	log.Debug().Msgf("Generating image #%d out of %d images", imageIndex, s.cfg.CollectionSize)

	startTime := time.Now().UTC()

	traits, err := s.traitService.GetRandomTraits()
	if err != nil {
		return errors.Annotate(err, "getting random traits failed")
	}

	if len(traits) == 0 {
		return errors.New("traits not found")
	}

	layers := make([]*model.ImageLayer, len(traits))
	for i, trait := range traits {
		layer, err := trait.ToImageLayer()
		if err != nil {
			return errors.Annotate(err, "converting image to layer failed")
		}

		layers[i] = layer
	}

	sortByPriority(layers)

	if err := s.generateImage(layers); err != nil {
		return errors.Annotate(err, "generating image failed")
	}

	log.Debug().Msgf("Finished generating image #%d (took %.3f seconds)", imageIndex, time.Since(startTime).Seconds())

	return nil
}

// private

func (s *generator) generateImage(layers []*model.ImageLayer) error {
	bgImg := image.NewRGBA(image.Rect(0, 0, s.cfg.ImageWidth, s.cfg.ImageHeight))

	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{C: color.Transparent}, image.Point{}, draw.Src)

	for _, layer := range layers {
		// Set the image offset.
		offset := image.Pt(layer.XPos, layer.YPos)

		// Combine the image.
		draw.Draw(bgImg, layer.Image.Bounds().Add(offset), layer.Image, image.Point{}, draw.Over)
	}

	buff := &bytes.Buffer{}
	if err := png.Encode(buff, bgImg); err != nil {
		return errors.Annotate(err, "encoding image to png failed")
	}

	fileName := filepath.Join(s.cfg.CollectionOutputDir, fmt.Sprintf("%d.png", imageIndex))
	if err := os.WriteFile(fileName, buff.Bytes(), 0777); err != nil {
		return errors.Annotate(err, "writing file failed")
	}

	return nil
}

func sortByPriority(list []*model.ImageLayer) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Priority < list[j].Priority
	})
}
