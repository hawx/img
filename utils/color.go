package utils

import (
	"image/color"
)

// NormalisedRGBA returns the RGBA colour channel values of a Color in a
// normalised form. The values are in non-premultiplied form and are in the
// range 0-255.
//
// This is modified from image/color.nrgbaModel
func NormalisedRGBA(c color.Color) (r, g, b, a uint32) {
	r, g, b, a = c.RGBA()

	r = r >> 8
	g = g >> 8
	b = b >> 8
	a = a >> 8

	if a == 0xff {
		return
	}

	if a == 0 {
		return 0, 0, 0, 0
	}

	// Since Color.RGBA returns a alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r = (r * 0xff) / a
	g = (g * 0xff) / a
	b = (b * 0xff) / a
	return
}

// NormalisedRGBAf returns the RGBA colour channel values as floating point
// numbers with values from 0 to 255.
func NormalisedRGBAf(c color.Color) (float64, float64, float64, float64) {
	r, g, b, a := NormalisedRGBA(c)
	return float64(r), float64(g), float64(b), float64(a)
}

// RatioRGBA returns the RGBA colour channel values as floating point numbers
// with values from 0 to 1.
func RatioRGBA(c color.Color) (float64, float64, float64, float64) {
	r, g, b, a := NormalisedRGBAf(c)
	return r / 255, g / 255, b / 255, a / 255
}

// Truncate takes a colour channel value and forces it into the range 0 to 255
// by setting any value below 0 to 0 and and any above 255 to 255.
func Truncate(n uint32) uint32 {
	if n < 0 {
		return 0
	} else if n > 255 {
		return 255
	}
	return n
}

// Truncatef is identical to Truncate but takes and returns a float64.
func Truncatef(n float64) float64 {
	if n < 0 {
		return 0
	} else if n > 255 {
		return 255
	}
	return n
}

// Closeness calculates the "closeness" of two colours by finding the sum of
// differences in each colour channel.
func Closeness(one, two color.Color) uint32 {
	a, b, c, d := NormalisedRGBA(one)
	w, x, y, z := NormalisedRGBA(two)

	return (a - w) + (b - x) + (c - y) + (d - z)
}

// Average takes a list of colours and returns the average. Given an empty list
// it returns Black.
func Average(cs ...color.Color) color.Color {
	if len(cs) < 1 {
		return color.Black
	}

	var red, green, blue, alpha uint32
	red = 0
	green = 0
	blue = 0
	alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := NormalisedRGBA(cs[i])

		red += r
		green += g
		blue += b
		alpha += a
	}

	return color.NRGBA{
		uint8(red / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}
