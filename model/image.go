package model

import (
	"image"
	"image/color"
)

type ImageLayer struct {
	Image    image.Image
	Priority int
	XPos     int
	YPos     int
}

type BgProperty struct {
	Width   int
	Length  int
	BgColor color.Color
}
