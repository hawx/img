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
	return utils.MapColor(img, func(c color.Color) color.Color {
		r, g, b, a := utils.NormalisedRGBA(c)
		grey := uint8(f(r, g, b))

		return color.NRGBA{grey, grey, grey, uint8(a)}
	})
}

// luminosityAlterer creates a function which scales the RGB channels by the
// ratios given, returning the value r*rM + g*gM + b*bM.
func luminosityAlterer(rM, gM, bM float64) pixelAlterer {
	return func(r, g, b uint32) uint32 {
		return uint32(float64(r) * rM + float64(g) * gM + float64(b) * bM)
	}
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

// Luminosity creates a greyscale version of the Image using the luminosity
// method. This uses a weighted average to account for human sensitivity to
// green above other colours.  Formula is R * 0.21 + G * 0.71 + B * 0.07.
func Luminosity(img image.Image) image.Image {
	return alterPixels(img, luminosityAlterer(0.21, 0.71, 0.07))
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

// Red creates a greyscale version of the Image using the values of the red
// channel.
func Red(img image.Image) image.Image {
	return alterPixels(img, luminosityAlterer(1, 0, 0))
}

// Green creates a greyscale version of the Image using the values of the green
// channel.
func Green(img image.Image) image.Image {
	return alterPixels(img, luminosityAlterer(0, 1, 0))
}

// Blue creates a greyscale version of the Image using the values of the blue
// channel.
func Blue(img image.Image) image.Image {
	return alterPixels(img, luminosityAlterer(0, 0, 1))
}

// Photoshop creates a greyscale version of the Image using the method
// (supposedly) used by Adobe Photoshop. It is simply a variation on the
// Luminosity method.
func Photoshop(img image.Image) image.Image {
	return alterPixels(img, luminosityAlterer(0.299, 0.587, 0.114))
}
