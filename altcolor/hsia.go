package altcolor

import (
	"image/color"
	"math"

	"hawx.me/code/img/utils"
)

// HSIA represents a colour with hue, saturation, intensity and alpha channels.
// All values are represented as floating point numbers with hue taking a value
// between 0 and 360, and the rest between 0 and 1.
//
// Sometimes called HSBA (Hue Saturation Brightness Alpha)
type HSIA struct {
	H, S, I, A float64
}

func (col HSIA) RGBA() (red, green, blue, alpha uint32) {
	var r, g, b float64
	h := col.H
	s := col.S
	i := col.I
	a := col.A

	// normalise h
	h = math.Mod(h, 360)

	// need h in radians for trig. calculations
	hrad := h

	// need hrad to be in interval [0, 2/3*Pi)
	if h >= 240 {
		hrad -= 240
	} else if h >= 120 {
		hrad -= 120
	}

	// finally do the conversion to radians
	hrad *= math.Pi / 180

	x := i * (1 - s)
	y := i * (1 + (s*math.Cos(hrad))/math.Cos(math.Pi/3-hrad))
	z := 3*i - (x + y)

	if h < 120 {
		b = x
		r = y
		g = z
	} else if h < 240 {
		r = x
		g = y
		b = z
	} else {
		g = x
		b = y
		r = z
	}

	red = uint32(utils.Truncatef(r*a*255)) << 8
	green = uint32(utils.Truncatef(g*a*255)) << 8
	blue = uint32(utils.Truncatef(b*a*255)) << 8
	alpha = uint32(a*255) << 8

	return
}

var HSIAModel color.Model = color.ModelFunc(hsiaModel)

func hsiaModel(c color.Color) color.Color {
	if _, ok := c.(HSIA); ok {
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

	// Work out intensity
	intensity := (r + g + b) / 3

	// Work out saturation
	saturation := 0.0
	if chroma != 0 {
		saturation = 1 - mini/intensity
	}

	// prefer positive hues
	if hue < 0 {
		hue += 360
	}

	return HSIA{hue, saturation, intensity, a}
}
