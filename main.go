package main

import (
	"os"
	"time"

	"github.com/Oxxyg33n/desert-island-go/configuration"
	"github.com/Oxxyg33n/desert-island-go/pinata"
	"github.com/Oxxyg33n/desert-island-go/repository"
	"github.com/Oxxyg33n/desert-island-go/service"
	"github.com/joho/godotenv"
	"github.com/juju/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	// Initialize dotenv
	if err := godotenv.Load(); err != nil {
		log.Fatal().
			Msg(errors.Annotate(err, "initializing dotenv failed").Error())
	}

	// Initialize configuration
	cfg := configuration.Configuration{}
	if err := cfg.New(); err != nil {
		log.Fatal().
			Msg(errors.Annotate(err, "initializing configuration failed").Error())
	}

	// Initialize repositories
	traitRepository := repository.NewTrait()
	groupRepository := repository.NewGroup()
	dnaRepository := repository.NewDNA()

	// Initialize services
	traitService := service.NewTrait(traitRepository, groupRepository, cfg.ImageWidth, cfg.ImageHeight)
	if err := traitService.Import(cfg.CollectionInputDir); err != nil {
		log.Fatal().
			Msg(errors.Annotate(err, "importing traits failed").Error())
	}

	dnaService := service.NewDNA(dnaRepository, cfg.CollectionDNAPrefix)
	metadataService := service.NewMetadata(cfg)
	generator := service.NewGenerator(cfg, traitService, dnaService, metadataService)
	ipfs := pinata.NewPinata(cfg.PinataAPIKey, cfg.PinataAPISecret, cfg.CollectionOutputDir)

	// Run generator
	startTime := time.Now().UTC()

	log.Info().
		Msgf("Generation started at %s", startTime.Format("2006-01-02 15:04:05"))
	log.Info().
		Msgf("Generating collection with size %d", cfg.CollectionSize)

	var totalImagesGenerated int
	for i := cfg.CollectionStartIndex; i <= cfg.CollectionSize; {
		if err := generator.Generate(i); err != nil {
			log.Error().
				Msg(errors.Annotate(err, "generating image failed").Error())

			continue
		}

		totalImagesGenerated++
		i++
	}

	log.Info().Msgf("Generation of %d images took %.3f second(-s)", totalImagesGenerated, time.Since(startTime).Seconds())

	if cfg.PinataUploadEnabled {
		startTime = time.Now().UTC()
		log.Info().
			Msgf("Uploading metadata to Pinata started at %s", startTime.Format("2006-01-02 15:04:05"))

		if err := ipfs.Upload(); err != nil {
			log.Fatal().
				Msg(errors.Annotate(err, "uploading metadata to Pinata failed").Error())
		}

		log.Info().
			Msgf("Finished uploading metadata to Pinata in %.3f seconds", time.Since(startTime).Seconds())
	}
}
