package utils

import (
	"image/color"
	"math"
)

// Converts RGBA color to HSLA,
//   h is in the range -pi..pi
//   s                 0..1
//   l                 0..1
//   a                 0..1
//
// This is not exact, uses formulas from
// http://www.quasimondo.com/archives/000696.php
func ToSimplifiedHSLA(c color.Color) (h, s, l, a float64) {
	r, g, b, a := RatioRGBA(c)

	l  =  r * 0.299 + g * 0.587 + b * 0.114
	u := -r * 0.1471376975169300226 - g * 0.2888623024830699774 + b * 0.436
	v :=  r * 0.615 - g * 0.514985734664764622 - b * 0.100014265335235378
	h  = math.Atan2(v, u)
	s  = math.Sqrt(u*u + v*v) * math.Sqrt(2)

	return
}

func ToSimplifiedRGBA(h, s, l, a float64) color.Color {
	u := math.Cos(h) * s
	v := math.Sin(h) * s
	r := l + 1.139837398373983740 * v
	g := l - 0.3946517043589703515 * u - 0.5805986066674976801 * v
	b := l + 2.03211091743119266 * u

	r  = Truncatef(r * 255)
	g  = Truncatef(g * 255)
	b  = Truncatef(b * 255)
	a  = a * 255

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Returns HSLA representation of given color with:
//   0 <= h < 360
//   0 <= l < 1
//   0 <= s < 1
//   0 <= a < 1
//
// Uses algorithm described on here
// http://en.wikipedia.org/wiki/HSL_and_HSV#General_approach
func ToHSLA(c color.Color) (h,s,l,a float64) {
	r, g, b, a := RatioRGBA(c)

	maxi   := Maxf(r, g, b)
	mini   := Minf(r, g, b)
	chroma := maxi - mini

	// Work out hue
	hdash  := 0.0
	if chroma == 0 {
		hdash = 0
	}else if maxi == r {
		hdash = math.Mod((g - b) / chroma, 6)
	} else if maxi == g {
		hdash = (b - r) / chroma + 2.0
	} else if maxi == b {
		hdash = (r - g) / chroma + 4.0
	}

	hue := hdash * 60

	if chroma == 0 {
		hue = 0
	}

	// Work out lightness
	lightness := 0.5 * (maxi + mini)

	// Work out saturation
	saturation := 0.0
	if chroma != 0 {
		saturation = chroma / (1 - math.Abs(2 * lightness - 1))
	}

	return hue, saturation, lightness, a
}

func ToRGBA(h,s,l,a float64) color.Color {
	h = math.Mod(h, 360)

	c := (1.0 - math.Abs(2.0 * l - 1.0)) * s

	hdash := h / 60.0
	x := c * (1 - math.Abs(math.Mod(hdash,2) - 1))

	var r, g, b float64
	if h == 0 {
		r = 0; g = 0; b = 0
	} else if hdash < 1 {
		r = c; g = x; b = 0
	} else if hdash < 2 {
		r = x; g = c; b = 0
	} else if hdash < 3 {
		r = 0; g = c; b = x
	} else if hdash < 4 {
		r = 0; g = x; b = c
	} else if hdash < 5 {
		r = x; g = 0; b = c
	} else if hdash < 6 {
		r = c; g = 0; b = x
	}

	m := l - 0.5 * c

	return color.RGBA{
		uint8(Truncatef((r + m) * 255)),
		uint8(Truncatef((g + m) * 255)),
		uint8(Truncatef((b + m) * 255)),
		uint8(a * 255),
	}
}
