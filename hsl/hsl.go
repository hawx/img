package hsl

import (
	"../utils"
	"image"
	"image/color"
)

// Shifts the hue of +img+ by +amount+. This should be a value between 0 and
// 360.
func Hue(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToHSLA(c)
		h += amount
		return utils.ToRGBA(h,s,l,a)
	}

	return utils.EachPixel(img, f)
}

// Adjusts the saturation of +img+ by +amount+.
func Saturation(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToHSLA(c)
		s += amount
		return utils.ToRGBA(h,s,l,a)
	}

	return utils.EachPixel(img, f)
}

// Adjusts the lightness of +img+ by +amount+. This produces slightly different
// results to brightness.Adjust.
func Lightness(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToHSLA(c)
		l += amount
		return utils.ToRGBA(h,s,l,a)
	}

	return utils.EachPixel(img, f)
}
