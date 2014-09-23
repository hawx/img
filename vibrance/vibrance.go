// Package vibrance provides functions to adjust the "vibrancy" of images.
package vibrance

import (
	"github.com/hawx/img/channel"
	"github.com/hawx/img/utils"
	"github.com/lucasb-eyer/go-colorful"

	"image"
	"image/color"
	"math"
)

// Adjust increases the saturation of the least saturated parts of the image
// while also reducing the lightness of these parts.
func Adjust(img image.Image, amount float64) image.Image {
	return utils.MapColor(img, AdjustC(amount))
}

// AdjustC returns a Composable function that increases the saturation and
// decreases the lightness of unsaturated colours.
//
// Uses the same method as Darktable:
// https://github.com/darktable-org/darktable/blob/24c4a087fd020df587b5260f438bfaf494203cec/src/iop/vibrance.c
func AdjustC(amount float64) utils.Composable {
	return func(c color.Color) color.Color {
		r, g, b, a := utils.RatioRGBA(c)

		ll, la, lb := colorful.Color{r, g, b}.Lab()

		// saturation weight [0, 1]
		sw := math.Sqrt(la*la+lb*lb) / 2
		ls := 1.0 - amount*sw*0.25
		ss := 1.0 + amount*sw

		ll *= ls
		la *= ss
		lb *= ss

		f := colorful.Lab(ll, la, lb)

		fr := utils.Truncatef(f.R * 255)
		fg := utils.Truncatef(f.G * 255)
		fb := utils.Truncatef(f.B * 255)

		return color.NRGBA{uint8(fr), uint8(fg), uint8(fb), uint8(a * 255)}
	}
}

// Exp increases the saturation of the least saturated parts of an image.
func Exp(img image.Image, amount float64) image.Image {
	return utils.MapColor(img, ExpC(amount))
}

// ExpC returns a Composable function that increases the saturation of
// unsaturated colours.
func ExpC(amount float64) utils.Composable {
	ch := channel.Saturation

	return func(c color.Color) color.Color {
		s := ch.Get(c)
		s = math.Pow(s, 1/amount)
		return ch.Set(c, s)
	}
}
