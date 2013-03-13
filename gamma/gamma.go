// Package gamma implements functions to adjust the gamma of an image.
package gamma

import (
	"github.com/hawx/img/utils"
	"github.com/hawx/img/greyscale"
	"image"
	"image/color"
	"math"
)

// Adjusts the gamma of the Image by the given value. A value less than 1.0
// darkens the image, whilst a gamma of greater than 1.0 lightens an image.
func Adjust(img image.Image, value float64) image.Image {
	return utils.MapColor(img, AdjustC(value))
}

func AdjustC(value float64) utils.Composable {
	return func(c color.Color) color.Color {
		r, g, b, a := utils.RatioRGBA(c)

		r = utils.Truncatef(math.Pow(r, 1/value) * 255)
		g = utils.Truncatef(math.Pow(g, 1/value) * 255)
		b = utils.Truncatef(math.Pow(b, 1/value) * 255)

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a * 255)}
	}
}

// Auto calculates the mean values of an image, then applies a gamma adjustment
// so that the mean colour in the image has a value of half.
func Auto(img image.Image) image.Image {
	luma := greyscale.Luminosity(img)
	greys := []uint32{}

	b := luma.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			g,_,_,_ := utils.NormalisedRGBA(luma.At(x, y))
			greys = append(greys, g)
		}
	}

	var averageGrey uint32 = 0
	for i := 0; i < len(greys); i++ {
		averageGrey += greys[i]
	}
	averageGrey /= uint32(len(greys))

	newGamma := 127.5 / float64(averageGrey)

	return Adjust(img, newGamma)
}
