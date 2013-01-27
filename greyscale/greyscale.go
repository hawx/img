// Package greyscale implements various functions to convert an image to
// greyscale.
package greyscale

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
)

type pixelAlterer func(r, g, b uint32) uint32

func alterColor(c color.Color, f pixelAlterer) color.Color {
	r,g,b,a := utils.NormalisedRGBA(c)
	grey := uint8(f(r, g, b))

	return color.NRGBA{grey, grey, grey, uint8(a)}
}

// luminosityAlterer creates a function which scales the RGB channels by the
// ratios given, returning the value r*rM + g*gM + b*bM.
func luminosityAlterer(rM, gM, bM float64) pixelAlterer {
	return func(r, g, b uint32) uint32 {
		return uint32(float64(r) * rM + float64(g) * gM + float64(b) * bM)
	}
}

func Mapper(f func(color.Color) color.Color) func(image.Image) image.Image {
	return func(img image.Image) image.Image {
		return utils.MapColor(img, f)
	}
}


// Average creates a greyscale version of the Image using the average method;
// it simply averages the RGB values.
var Average = Mapper(AverageC)

func AverageC(c color.Color) color.Color {
	return alterColor(c, func(r,g,b uint32) uint32 {
		return (r + g + b) / 3
	})
}

// Lightness creates a grayscale version of the Image using the lightness
// method. That is, it averages the most prominent and least prominent colours.
var Lightness = Mapper(LightnessC)

func LightnessC(c color.Color) color.Color {
	return alterColor(c, func(r,g,b uint32) uint32 {
		return (utils.Max(r, g, b) + utils.Min(r, g, b)) / 2
	})
}

// Luminosity creates a greyscale version of the Image using the luminosity
// method. This uses a weighted average to account for human sensitivity to
// green above other colours.  Formula is R * 0.21 + G * 0.71 + B * 0.07.
var Luminosity = Mapper(LuminosityC)

func LuminosityC(c color.Color) color.Color {
	return alterColor(c, luminosityAlterer(0.21, 0.71, 0.07))
}

// Maximal creates a greyscale version of the Image by taking the maximum of the
// RGB channels and using that value for each channel. Produces a lighter
// image.
var Maximal = Mapper(MaximalC)

func MaximalC(c color.Color) color.Color {
	return alterColor(c, func(r, g, b uint32) uint32 {
		return utils.Max(r, g, b)
	})
}

// Minimal creates a greyscale version of the Image by taking the minimum value
// of the RGB channels and using that for each channel. Produces a darker
// image.
var Minimal = Mapper(MinimalC)

func MinimalC(c color.Color) color.Color {
	return alterColor(c, func(r, g, b uint32) uint32 {
		return utils.Min(r, g, b)
	})
}

// Red creates a greyscale version of the Image using the values of the red
// channel.
var Red = Mapper(RedC)

func RedC(c color.Color) color.Color {
	return alterColor(c, luminosityAlterer(1, 0, 0))
}

// Green creates a greyscale version of the Image using the values of the green
// channel.
var Green = Mapper(GreenC)

func GreenC(c color.Color) color.Color {
	return alterColor(c, luminosityAlterer(0, 1, 0))
}

// Blue creates a greyscale version of the Image using the values of the blue
// channel.
var Blue = Mapper(BlueC)

func BlueC(c color.Color) color.Color {
	return alterColor(c, luminosityAlterer(0, 0, 1))
}

// Photoshop creates a greyscale version of the Image using the method
// (supposedly) used by Adobe Photoshop. It is simply a variation on the
// Luminosity method.
var Photoshop = Mapper(PhotoshopC)

func PhotoshopC(c color.Color) color.Color {
	return alterColor(c, luminosityAlterer(0.299, 0.587, 0.114))
}
