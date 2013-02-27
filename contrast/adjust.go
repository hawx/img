// Package contrast implements functions to adjust the contrast of an image.
package contrast

import (
	"github.com/hawx/img/altcolor"
	"github.com/hawx/img/utils"
	"image"
	"image/color"
	"math"
)

const Epsilon = 1.0e-10

// Adjust changes the contrast in the Image. A value of 0 has no effect.
func Adjust(img image.Image, value float64) image.Image {
	return utils.MapColor(img, AdjustC(value))
}

func AdjustC(value float64) utils.Composable {
	return func(c color.Color) color.Color {
		d := altcolor.HSIAModel.Convert(c).(altcolor.HSIA)

		d.I += 0.5 * value * (0.5 * (math.Sin(math.Pi * (d.I - 0.5)) + 1.0) - d.I)
		if d.I > 1.0 { d.I = 1.0 } else if d.I < 0.0 { d.I = 0.0 }

		return d
	}
}


// Linear adjusts the contrast using a linear function. A value of 1 has no
// effect, and a value of 0 will return a grey image.
func Linear(img image.Image, value float64) image.Image {
	return utils.MapColor(img, LinearC(value))
}

func LinearC(value float64) utils.Composable {
	return func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)

		r = utils.Truncatef((((r - 0.5) * value) + 0.5) * 255)
		g = utils.Truncatef((((g - 0.5) * value) + 0.5) * 255)
		b = utils.Truncatef((((b - 0.5) * value) + 0.5) * 255)
		a = a * 255

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
}


// Sigmoidal adjusts the contrast in a non-linear way. Factor sets how much to
// increase the contrast, midpoint sets where midtones fall in the resultant
// image.
func Sigmoidal(img image.Image, factor, midpoint float64) image.Image {
	return utils.MapColor(img, SigmoidalC(factor, midpoint))
}

func SigmoidalC(factor, midpoint float64) utils.Composable {
	sigmoidal := func(x float64) float64 {
		return 1.0 / (1.0 + math.Exp(factor * (midpoint - x)))
	}

	// Pre-compute useful terms
	sig0 := sigmoidal(0.0)
	sig1 := sigmoidal(1.0)

	var scaledSigmoidal func(float64) float64

	if factor == 0 {
		scaledSigmoidal = func(x float64) float64 {
			return x
		}

	} else if factor > 0 {
		scaledSigmoidal = func(x float64) float64 {
			return (sigmoidal(x) - sig0) / (sig1 - sig0)
		}

	} else {
		scaledSigmoidal = func(x float64) float64 {
			argument := (sig1 - sig0) * x + sig0
			var clamped float64
			if argument < Epsilon {
				clamped = Epsilon
			} else {
				if argument > 1 - Epsilon {
					clamped = 1 - Epsilon
				} else {
					clamped = argument
				}
			}

			return midpoint - math.Log(1.0 / clamped - 1.0) / factor
		}
	}

	return func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)

		r = utils.Truncatef(scaledSigmoidal(r) * 255)
		g = utils.Truncatef(scaledSigmoidal(g) * 255)
		b = utils.Truncatef(scaledSigmoidal(b) * 255)
		a = a * 255

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
}
