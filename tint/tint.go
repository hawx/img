// Package tint implements a single function to tint an image with a specified
// colour.
package tint

import (
	"github.com/hawx/img/blend"
	"image"
	"image/color"
)

func Tint(img image.Image, with color.Color) image.Image {
	blendLayer := image.NewUniform(with)
	return blend.Normal(img, blendLayer)
}
