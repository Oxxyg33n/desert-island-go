package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Oxxyg33n/desert-island-go/model"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	var ipfsLink string
	flag.StringVar(&ipfsLink, "link", "", "Pass in IPFS link in format of ipfs://<link>")

	dirEntry, err := os.ReadDir("../../output/traits")
	if err != nil {
		panic(err)
	}

	for _, f := range dirEntry {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}

		if err := processFile(f.Name(), ipfsLink); err != nil {
			panic(errors.Annotate(err, "processing file failed"))
		}
	}
}

// private

func processFile(fileName, ipfsLink string) error {
	fileNameWithPath := "../../output/traits/" + fileName
	log.Debug().Msgf("opening file %s", fileNameWithPath)

	fBytes, err := ioutil.ReadFile(fileNameWithPath)
	if err != nil {
		log.Error().Msg(errors.Annotate(err, "reading file failed").Error())
		return err
	}

	fmt.Println("here")
	fmt.Print(string(fBytes))

	var metadata model.ERCMetadata
	if err := json.Unmarshal(fBytes, &metadata); err != nil {
		log.Error().Msg(errors.Annotate(err, "unmarshalling json failed").Error())
		return err
	}

	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	metadata.Image = fmt.Sprintf("%s/%s.png", ipfsLink, fileNameWithoutExt)

	b, err := json.Marshal(metadata)
	if err != nil {
		log.Error().Msg(errors.Annotate(err, "marshalling json failed").Error())
		return err
	}

	if err := os.WriteFile(fileNameWithPath, b, 0666); err != nil {
		log.Error().Msg(errors.Annotate(err, "writing file bytes failed").Error())
		return err
	}

	return nil
}
