package model

type ERCMetadata struct {
	Image      string     `json:"image"`
	Attributes []ERCTrait `json:"trait"`
}
