package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Oxxyg33n/desert-island-go/model"
	pinatamodel "github.com/Oxxyg33n/desert-island-go/pinata/model"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

type IPFS interface {
	Upload() error
}

var _ IPFS = &pinata{}

type pinata struct {
	client    http.Client
	apiKey    string
	apiSecret string
	inputDir  string
}

func NewPinata(apiKey, apiSecret, inputDir string) IPFS {
	return &pinata{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		inputDir:  inputDir,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *pinata) Upload() error {
	var files []*file
	if err := filepath.Walk(s.inputDir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		if filepath.Ext(fileName) != ".png" {
			return nil
		}

		files = append(files, &file{
			path: path,
			name: fileName,
		})

		return nil
	}); err != nil {
		return errors.Annotate(err, "walking directory failed")
	}

	log.Debug().Msgf("Prepared %d files to upload", len(files))

	for _, f := range files {
		if err := s.uploadImage(f); err != nil {
			return errors.Annotate(err, "uploading image failed")
		}
	}

	return nil
}

// private

type file struct {
	path string
	name string
}

func (s *pinata) uploadImage(f *file) error {
	key := strings.TrimSuffix(f.name, filepath.Ext(f.name))

	b, err := ioutil.ReadFile(f.path)
	if err != nil {
		return errors.Annotate(err, "reading image file failed")
	}

	ipfsImageHash, err := s.pinFile(f.name, b)
	if err != nil {
		return errors.Annotate(err, "pinning file failed")
	}

	b, err = ioutil.ReadFile(strings.Replace(f.path, ".png", ".json", 1))
	if err != nil {
		return errors.Annotate(err, "reading traits file failed")
	}

	var ercTraits []model.ERCTrait
	if err := json.Unmarshal(b, &ercTraits); err != nil {
		return errors.Annotate(err, "unmarshalling json failed")
	}

	erc721Metadata := &model.ERCMetadata{
		Image:      fmt.Sprintf("ipfs://%s", ipfsImageHash),
		Attributes: ercTraits,
	}

	erc721MetadataBytes, err := json.Marshal(erc721Metadata)
	if err != nil {
		return errors.Annotate(err, "")
	}

	filename := fmt.Sprintf("%s/%s.meta.json", s.inputDir, key)
	if err := os.WriteFile(filename, erc721MetadataBytes, 0777); err != nil {
		return errors.Annotate(err, "writing file failed")
	}

	log.Debug().
		Msgf("Uploaded image %s to IPFS", f.name)

	return nil
}

const pinFileURL = "https://api.pinata.cloud/pinning/pinFileToIPFS"

func (s *pinata) pinFile(fileName string, data []byte) (string, error) {
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)

	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return "", errors.Annotate(err, "creating form file failed")
	}

	if _, err := fileWriter.Write(data); err != nil {
		return "", errors.Annotate(err, "writing data failed")
	}

	fileWriter, err = bodyWriter.CreateFormField("pinataOptions")
	if err != nil {
		return "", errors.Annotate(err, "creating form field failed")
	}

	pinataOptions := pinatamodel.Options{
		WrapWithDirectory: false,
		CIDVersion:        pinatamodel.CIDVersion1,
	}

	b, err := json.Marshal(pinataOptions)
	if err != nil {
		return "", errors.Annotate(err, "marshalling json failed")
	}

	if _, err := fileWriter.Write(b); err != nil {
		return "", errors.Annotate(err, "writing field failed")
	}

	contentType := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		return "", errors.Annotate(err, "closing writer failed")
	}

	req, err := http.NewRequest(http.MethodPost, pinFileURL, buf)
	if err != nil {
		return "", errors.Annotate(err, "creating new POST http request failed")
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("pinata_api_key", s.apiKey)
	req.Header.Set("pinata_secret_api_key", s.apiSecret)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", errors.Annotate(err, "doing request failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Annotate(err, "reading response body failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("http request failed with status code %d: message=%s", resp.StatusCode, string(body))
	}

	var pinataResp pinatamodel.PinataResponse
	if err = json.NewDecoder(bytes.NewReader(body)).Decode(&pinataResp); err != nil {
		return "", errors.Annotate(err, "unmarshalling json failed")
	}

	return pinataResp.IPFSHash, nil
}
