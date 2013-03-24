package utils

import (
	"image/color"
)

// NormalisedRGBA returns the RGBA colour channel values of a Color in a
// normalised form. The values are in non-premultiplied form and are in the
// range 0-255.
func NormalisedRGBA(c color.Color) (rn, gn, bn, an uint32) {
	d := color.NRGBAModel.Convert(c).(color.NRGBA)
	r := d.R; g := d.G; b := d.B; a := d.A

	// Need to do some crazy type conversions first
	rn = uint32(uint8(r))
	gn = uint32(uint8(g))
	bn = uint32(uint8(b))
	an = uint32(uint8(a))

	return
}

// NormalisedRGBAf returns the RGBA colour channel values as floating point
// numbers with values from 0 to 255.
func NormalisedRGBAf(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := NormalisedRGBA(c)

	rn = float64(r)
	gn = float64(g)
	bn = float64(b)
	an = float64(a)

	return
}

// RatioRGBA returns the RGBA colour channel values as floating point numbers
// with values from 0 to 1.
func RatioRGBA(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := NormalisedRGBAf(c)

	rn = r / 255
	gn = g / 255
	bn = b / 255
	an = a / 255

	return
}

// Truncate takes a colour channel value and forces it into the range 0 to 255
// by setting any value below 0 to 0 and and any above 255 to 255.
func Truncate(n uint32) uint32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

// Truanctef is identical to Truncate but takes and returns a float64.
func Truncatef(n float64) float64 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
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
func Average(cs... color.Color) color.Color {
	if len(cs) < 1 {
		return color.Black
	}

	var red, green, blue, alpha uint32
	red = 0; green = 0; blue = 0; alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := NormalisedRGBA(cs[i])

		red += r; green += g; blue += b; alpha += a
	}

	return color.NRGBA{
		uint8(red   / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue  / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}
