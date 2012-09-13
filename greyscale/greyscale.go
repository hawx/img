// Package greyscale implements various functions to convert an image to
// greyscale.
package greyscale

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
)

type pixelAlterer func(r, g, b uint32) uint32

func alterPixels(img image.Image, f pixelAlterer) image.Image {
	return utils.EachPixel(img, func(c color.Color) color.Color {
		r, g, b, a := utils.NormalisedRGBA(c)
		grey := uint8(f(r, g, b))

		return color.NRGBA{grey, grey, grey, uint8(a)}
	})
}


// Average creates a greyscale version of the Image using the average method;
// it simply averages the RGB values.
func Average(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return (r + g + b) / 3
	})
}

// Lightness creates a grayscale version of the Image using the lightness
// method. That is, it averages the most prominent and least prominent colours.
func Lightness(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return (utils.Max(r, g, b) + utils.Min(r, g, b)) / 2
	})
}

// Lightness creates a greyscale version of the Image using the luminosity
// method. This uses a weighted average to account for human sensitivity to
// green above other colours.  Formula is R * 0.21 + G * 0.71 + B * 0.07.
func Luminosity(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return uint32(float64(r) * 0.21 + float64(g) * 0.71 + float64(b) * 0.07)
	})
}

// Maximal creates a greyscale version of the Image by taking the maximum of the
// RGB channels and using that value for each channel. Produces a lighter
// image.
func Maximal(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return utils.Max(r, g, b)
	})
}

// Minimal creates a greyscale version of the Image by taking the minimum value
// of the RGB channels and using that for each channel. Produces a darker
// image.
func Minimal(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return utils.Min(r, g, b)
	})
}

// Photoshop creates a greyscale version of the Image using the method
// (supposedly) used by Adobe Photoshop. It is simply a variation on the
// Luminosity method.
func Photoshop(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return uint32(float64(r) * 0.299 + float64(g) * 0.587 + float64(b) * 0.114)
	})
}
