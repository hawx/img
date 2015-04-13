// Package levels provides functions to alter the distibution of values of a
// color channel in an image.
package levels

import (
	"github.com/hawx/img/channel"
	"github.com/hawx/img/utils"
	"image"
	"image/color"
)

// linearScale scales the value given so that the range min to max is scaled to
// 0 to 1.
func linearScale(value, min, max float64) float64 {
	return (value - min) * (1 / (max - min))
}

// http://en.wikipedia.org/wiki/Histogram_equalization
func Equalise(img image.Image) image.Image {
	return img
}

func Auto(img image.Image, ch channel.Channel) image.Image {
	var lightest, darkest float64
	lightest = 0.0
	darkest = 1.0

	utils.PEachColor(img, func(c color.Color) {
		v := ch.Get(c)

		if v > lightest {
			lightest = v
		}
		if v < darkest {
			darkest = v
		}
	})

	// Use linear stretching algorithm
	//   v = (v - inLow) * ((outUp - outLow) / (inUp - inLow)) + outLow
	return utils.MapColor(img, func(c color.Color) color.Color {
		v := ch.Get(c)
		v = linearScale(v, darkest, lightest)

		return ch.Set(c, v)
	})
}

func AutoWhite(img image.Image, ch channel.Channel) image.Image {
	var lightest float64 = 0.0

	utils.PEachColor(img, func(c color.Color) {
		v := ch.Get(c)

		if v > lightest {
			lightest = v
		}
	})

	return SetWhite(img, ch, lightest)
}

// AutoBlack finds the darkest colour in the image and makes it black, adjusting
// the colours of every other point to achieve the same distribution.
func AutoBlack(img image.Image, ch channel.Channel) image.Image {
	var darkest float64 = 1.0

	utils.PEachColor(img, func(c color.Color) {
		v := ch.Get(c)

		if v < darkest {
			darkest = v
		}
	})

	return SetBlack(img, ch, darkest)
}

func SetBlack(img image.Image, ch channel.Channel, darkest float64) image.Image {
	return utils.MapColor(img, SetBlackC(ch, darkest))
}

func SetBlackC(ch channel.Channel, darkest float64) utils.Composable {
	return func(c color.Color) color.Color {
		v := ch.Get(c)
		v = linearScale(v, darkest, 1)

		return ch.Set(c, v)
	}
}

func SetWhite(img image.Image, ch channel.Channel, lightest float64) image.Image {
	return utils.MapColor(img, SetWhiteC(ch, lightest))
}

func SetWhiteC(ch channel.Channel, lightest float64) utils.Composable {
	return func(c color.Color) color.Color {
		v := ch.Get(c)
		v = linearScale(v, 0, lightest)

		return ch.Set(c, v)
	}
}

func SetCurve(img image.Image, ch channel.Channel, curve *Curve) image.Image {
	return utils.MapColor(img, SetCurveC(ch, curve))
}

func SetCurveC(ch channel.Channel, curve *Curve) utils.Composable {
	return func(c color.Color) color.Color {
		v := ch.Get(c)
		v = curve.Value(v)

		return ch.Set(c, v)
	}
}
