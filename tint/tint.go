// Package tint implements a single function to tint an image with a specified
// colour.
package tint

import (
	"image"
	"image/color"

	"github.com/hawx/img/blend"
)

// Tint adds a colored tint to the image.
func Tint(img image.Image, with color.Color) image.Image {
	blendLayer := image.NewUniform(with)
	return blend.LinearLight(img, blendLayer)
}
