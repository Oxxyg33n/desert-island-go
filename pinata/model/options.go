package model

type Options struct {
	WrapWithDirectory bool       `json:"wrapWithDirectory"`
	CIDVersion        CIDVersion `json:"cidVersion"`
}

type CIDVersion string

const (
	CIDVersion0 CIDVersion = "0"
	CIDVersion1 CIDVersion = "1"
)
