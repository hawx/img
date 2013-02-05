package channel

import (
	"github.com/hawx/img/utils"
	"github.com/hawx/img/altcolor"
	"image/color"
	"math"
)

// Red applies the Adjuster to the red channel of each pixel in the Image.
var Red = utils.MapAdjuster(RedC)

func RedC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		r = adj(r)
		if r > 1 { r = 1 } else if r < 0 { r = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}
}

// Green applies the Adjuster to the green channel of each pixel in the Image.
var Green = utils.MapAdjuster(GreenC)

func GreenC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		g = adj(g)
		if g > 1 { g = 1 } else if g < 0 { g = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}
}

// Blue applies the Adjuster to the blue channel of each pixel in the Image.
var Blue = utils.MapAdjuster(BlueC)

func BlueC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		b = adj(b)
		if b > 1 { b = 1 } else if b < 0 { b = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}
}

// Alpha applies the Adjuster to the alpha channel of each pixel in the Image.
var Alpha = utils.MapAdjuster(AlphaC)

func AlphaC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		a = adj(a)
		if a > 1 { a = 1 } else if a < 0 { a = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}
}

// Hue shifts the hue of the Image using the function given.
var Hue = utils.MapAdjuster(HueC)

func HueC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
		h.H = math.Mod(adj(h.H), 360)
		return h
	}
}

// Saturation adjusts the saturation of the Image using the function given.
var Saturation = utils.MapAdjuster(SaturationC)

func SaturationC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
		h.S = adj(h.S)
		if h.S > 1 { h.S = 1 }
		if h.S < 0 { h.S = 0 }
		return h
	}
}

// Lightness adjusts the lightness of the Image using the function given.
var Lightness = utils.MapAdjuster(LightnessC)

func LightnessC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		h := hsla.HSLAModel.Convert(c).(hsla.HSLA)
		h.L = adj(h.L)
		if h.L > 1 { h.L = 1 }
		if h.L < 0 { h.L = 0 }
		return h
	}
}
