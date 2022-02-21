package service

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/Oxxyg33n/desert-island-go/configuration"
	"github.com/juju/errors"
)

type ImageGenerator interface {
	Generate(imageIndex int) error
}

var _ ImageGenerator = &generator{}

type generator struct {
	cfg configuration.Configuration
}

func NewGenerator(cfg configuration.Configuration) ImageGenerator {
	return &generator{
		cfg: cfg,
	}
}

func (g *generator) Generate(imageIndex int) error {
	bgImg := image.NewRGBA(image.Rect(0, 0, g.cfg.ImageWidth, g.cfg.ImageHeight))

	buff := &bytes.Buffer{}
	if err := png.Encode(buff, bgImg); err != nil {
		return errors.Annotate(err, "encoding image to png failed")
	}

	fileName := filepath.Join("output", fmt.Sprintf("%d.png", imageIndex))
	if err := os.WriteFile(fileName, buff.Bytes(), 0777); err != nil {
		return errors.Annotate(err, "writing file failed")
	}

	return nil
}
