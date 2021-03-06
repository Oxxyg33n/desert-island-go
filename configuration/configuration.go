package configuration

import (
	"github.com/caarlos0/env/v6"
	"github.com/juju/errors"
)

type Configuration struct {
	CollectionSize        int    `env:"COLLECTION_SIZE" envDefault:"10"`
	CollectionDNAPrefix   string `env:"COLLECTION_DNA_PREFIX" envDefault:"0"`
	CollectionName        string `env:"COLLECTION_NAME"`
	CollectionDescription string `env:"COLLECTION_DESCRIPTION"`
	CollectionStartIndex  int    `env:"COLLECTION_START_INDEX" envDefault:"0"`
	CollectionOutputDir   string `env:"COLLECTION_OUTPUT_DIR" envDefault:"output"`
	CollectionInputDir    string `env:"COLLECTION_INPUT_DIR" envDefault:"input"`

	ImageWidth  int `env:"IMAGE_WIDTH" envDefault:"1000"`
	ImageHeight int `env:"IMAGE_HEIGHT" envDefault:"1000"`

	PinataUploadEnabled bool   `env:"PINATA_UPLOAD_ENABLED" envDefault:"false"`
	PinataAPISecret     string `env:"PINATA_API_SECRET"`
	PinataAPIKey        string `env:"PINATA_API_KEY"`
}

func (c *Configuration) New() error {
	if err := env.Parse(c); err != nil {
		return errors.Annotate(err, "parsing environment configuration failed")
	}

	return nil
}
