package model

import "fmt"

type ERCMetadata struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	Attributes  []ERCTrait `json:"attributes"`
}

func NewERCMetadata(
	collectionName, collectionDescription, ipfsImageCidr string,
	imageIndex int,
	attributes []ERCTrait,
) *ERCMetadata {
	return &ERCMetadata{
		Name:        fmt.Sprintf("%s - #%d", collectionName, imageIndex),
		Description: collectionDescription,
		Image:       ipfsImageCidr,
		Attributes:  attributes,
	}
}
