package utils

import (
	"os"
	"math"
	"image"
	"image/png"
	"image/color"
)

func ReadStdin() image.Image {
	img, _ := png.Decode(os.Stdin)

	return img
}

func WriteStdout(img image.Image) {
	png.Encode(os.Stdout, img)
}

func NormalisedRGBA(c color.Color) (rn, gn, bn, an uint32) {
	r, g, b, a := c.RGBA()

	// Need to do some crazy type conversions first
	rn = uint32(uint8(r))
	gn = uint32(uint8(g))
	bn = uint32(uint8(b))
	an = uint32(uint8(a))

	return
}

func RatioRGBA(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := c.RGBA()

	rn = float64(uint8(r)) / 255
	gn = float64(uint8(g)) / 255
	bn = float64(uint8(b)) / 255
	an = float64(uint8(a)) / 255

	return
}

// Converts RGBA color to HSLA,
//   h is in the range -pi..pi
//   s                 0..1
//   l                 0..1
//   a                 0..1
//
// This is not exact, uses formulas from
// http://www.quasimondo.com/archives/000696.php
// will update to proper method later.
func ToHSLA(c color.Color) (h, s, l, a float64) {
	r, g, b, a := RatioRGBA(c)

	l  =  r * 0.299 + g * 0.587 + b * 0.114
	u := -r * 0.1471376975169300226 - g * 0.2888623024830699774 + b * 0.436
	v :=  r * 0.615 - g * 0.514985734664764622 - b * 0.100014265335235378
	h  = math.Atan2(v, u)
	s  = math.Sqrt(u*u + v*v) * math.Sqrt(2)

	return
}

func ToRGBA(h, s, l, a float64) color.Color {
	u := math.Cos(h) * s
	v := math.Sin(h) * s
	r := l + 1.139837398373983740 * v
	g := l - 0.3946517043589703515 * u - 0.5805986066674976801 * v
	b := l + 2.03211091743119266 * u

	r  = TruncateFloat(r * 255)
	g  = TruncateFloat(g * 255)
	b  = TruncateFloat(b * 255)
	a  = a * 255

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func TruncateInt(n uint32) uint32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

func TruncateFloat(n float64) float64 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

func Closeness(one, two color.Color) uint32 {
	a, b, c, d := NormalisedRGBA(one)
	w, x, y, z := NormalisedRGBA(two)

	return (a - w) + (b - x) + (c - y) + (d - z)
}

func Average(cs []color.Color) color.Color {
	var red, green, blue, alpha uint32
	red = 0; green = 0; blue = 0; alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := NormalisedRGBA(cs[i])

		red += r; green += g; blue += b; alpha += a
	}

	return color.RGBA{
		uint8(red   / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue  / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}

func Min(ns... uint32) (n uint32) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] < n {
			n = ns[i]
		}
	}
	return
}

func Max(ns... uint32) (n uint32) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] > n {
			n = ns[i]
		}
	}
	return
}
