// Package gamma implements a signle function to adjust the gamma of an image.
package gamma

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
	"math"
)

// Adjusts the gamma of the Image by the given value. A value less than 1.0
// darkens the image, whilst a gamma of greater than 1.0 lightens an image.
func Adjust(img image.Image, value float64) image.Image {
	f := func(c color.Color) color.Color {
		r, g, b, a := utils.RatioRGBA(c)

		r = utils.Truncatef(math.Pow(r, 1/value) * 255)
		g = utils.Truncatef(math.Pow(g, 1/value) * 255)
		b = utils.Truncatef(math.Pow(b, 1/value) * 255)

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a * 255)}
	}

	return utils.EachPixel(img, f)
}
