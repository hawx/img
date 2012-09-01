package greyscale

import (
	"../utils"
	"image"
	"image/color"
)

type pixelAlterer func(r, g, b uint32) uint32

func alterPixels(img image.Image, f pixelAlterer) image.Image {
	return utils.EachPixel(img, func(c color.Color) color.Color {
		r, g, b, a := utils.NormalisedRGBA(c)
		grey := uint8(f(r, g, b))

		return color.RGBA{grey, grey, grey, uint8(a)}
	})
}


// Creates a greyscale version of +img+ using the average method, simply averages
// R, G and B values.
func Average(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return (r + g + b) / 3
	})
}

// Creates a grayscale version of +img+ using the lightness method. That is, it
// averages the most prominent and least prominent colours.
func Lightness(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return (utils.Max(r, g, b) + utils.Min(r, g, b)) / 2
	})
}

// Creates a greyscale version of +img+ using the luminosity method. This uses a
// weighted average to account for humans sensitivity to green above other colours.
// Formula is R * 0.21 + G * 0.71 + B * 0.07.
func Luminosity(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return uint32(float64(r) * 0.21 + float64(g) * 0.71 + float64(b) * 0.07)
	})
}

// Takes the maximum of the r,g,b channels and uses that value for each channel.
// Produces a lighter image.
func Maximal(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return utils.Max(r, g, b)
	})
}

// Takes the minimum value of the r,g,b channels and uses that for each channel.
// Produces a darker image.
func Minimal(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return utils.Min(r, g, b)
	})
}

// Supposed photoshop luminosity method for greyscale.
func Photoshop(img image.Image) image.Image {
	return alterPixels(img, func(r, g, b uint32) uint32 {
		return uint32(float64(r) * 0.299 + float64(g) * 0.587 + float64(b) * 0.114)
	})
}
