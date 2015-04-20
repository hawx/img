package altcolor

import (
	"image/color"
	"math"

	"hawx.me/code/img/utils"
)

// HSVA represents a colour with hue, saturation, value and alpha channels.
// All values are represented as floating point numbers with hue taking a value
// between 0 and 360, and the rest between 0 and 1.
type HSVA struct {
	H, S, V, A float64
}

func (col HSVA) RGBA() (red, green, blue, alpha uint32) {
	h := col.H
	s := col.S
	v := col.V
	a := col.A

	h = math.Mod(h, 360)

	c := v * s

	hdash := h / 60.0
	x := c * (1 - math.Abs(math.Mod(hdash, 2)-1))

	var r, g, b float64
	if hdash < 1 {
		r = c
		g = x
		b = 0
	} else if hdash < 2 {
		r = x
		g = c
		b = 0
	} else if hdash < 3 {
		r = 0
		g = c
		b = x
	} else if hdash < 4 {
		r = 0
		g = x
		b = c
	} else if hdash < 5 {
		r = x
		g = 0
		b = c
	} else if hdash < 6 {
		r = c
		g = 0
		b = x
	}

	m := v - c

	red = uint32(utils.Truncatef((r+m)*a*255)) << 8
	green = uint32(utils.Truncatef((g+m)*a*255)) << 8
	blue = uint32(utils.Truncatef((b+m)*a*255)) << 8
	alpha = uint32(a*255) << 8

	return
}

var HSVAModel color.Model = color.ModelFunc(hsvaModel)

func hsvaModel(c color.Color) color.Color {
	if _, ok := c.(HSVA); ok {
		return c
	}

	r, g, b, a := utils.RatioRGBA(c)

	maxi := utils.Maxf(r, g, b)
	mini := utils.Minf(r, g, b)
	chroma := maxi - mini

	// Work out hue
	hdash := 0.0
	if chroma == 0 {
		hdash = 0
	} else if maxi == r {
		hdash = math.Mod((g-b)/chroma, 6)
	} else if maxi == g {
		hdash = (b-r)/chroma + 2.0
	} else if maxi == b {
		hdash = (r-g)/chroma + 4.0
	}

	hue := hdash * 60

	if chroma == 0 {
		hue = 0
	}

	// Work out value
	value := maxi

	// Work out saturation
	saturation := 0.0
	if chroma != 0 {
		saturation = chroma / value
	}

	return HSVA{hue, saturation, value, a}
}
