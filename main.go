package main

import (
	"os"
	"time"

	"github.com/Oxxyg33n/desert-island-go/configuration"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal().
			Msg(errors.Annotate(err, "initializing dotenv failed").Error())
	}

	cfg := configuration.Configuration{}
	if err := cfg.New(); err != nil {
		log.Fatal().
			Msg(errors.Annotate(err, "initializing configuration failed").Error())
	}

	traitService := service.NewTrait()
	if err := traitService.Import(cfg.InputDir); err != nil {
		log.Fatal().
			Msg(errors.Annotate(err, "importing traits failed").Error())
	}

	startTime := time.Now().UTC()

	log.Info().
		Msgf("Generation started at %s", startTime.Format("2005-01-02 15:04:05"))
	log.Info().
		Msgf("Generating collection with size %d", cfg.CollectionSize)

	generator := service.NewGenerator(cfg)
	for i := 0; i < cfg.CollectionSize; i++ {
		if err := generator.Generate(i); err != nil {
			log.Fatal().
				Msg(errors.Annotate(err, "generating image failed").Error())
		}

	}

	log.Info().Msgf("Generation took %d second(-s)", time.Since(startTime).Milliseconds())
}
