// Package channel provides functions to alter the values of separate color
// channels in an image.
package channel

import (
	"image"
	"image/color"
	"math"

	"github.com/hawx/img/altcolor"
	"github.com/hawx/img/utils"
)

// Adjust applies the given Adjuster to the Image on only the Channels specified.
func Adjust(img image.Image, adj utils.Adjuster, chs ...Channel) image.Image {
	return utils.MapColor(img, AdjustC(adj, chs...))
}

// AdjustC returns a Composable function that applies the given Adjuster to the
// Channels.
func AdjustC(adj utils.Adjuster, chs ...Channel) utils.Composable {
	return func(c color.Color) color.Color {
		for _, ch := range chs {
			v := ch.Get(c)
			c = ch.Set(c, adj(v))
		}

		return c
	}
}

// A Channel provides Get and Set methods, for accessing and modifying the value
// of a Color in the corresponding color channel. Get and Set both work on
// values in the range [0,1]; that is, the value of the channel will be
// scaled to be in that range.
type Channel interface {
	Get(color.Color) float64
	Set(color.Color, float64) color.Color
}

var (
	Red        = redCh{}
	Green      = greenCh{}
	Blue       = blueCh{}
	Alpha      = alphaCh{}
	Hue        = hueCh{}
	Saturation = saturationCh{}
	Lightness  = lightnessCh{}
	Intensity  = intensityCh{}

	// Alias
	Brightness = Intensity
)

type redCh struct{}

func (_ redCh) Get(c color.Color) float64 {
	r, _, _, _ := utils.RatioRGBA(c)
	return r
}

func (_ redCh) Set(c color.Color, v float64) color.Color {
	_, g, b, a := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(v), uint8(g), uint8(b), uint8(a)}
}

type greenCh struct{}

func (_ greenCh) Get(c color.Color) float64 {
	_, g, _, _ := utils.RatioRGBA(c)
	return g
}

func (_ greenCh) Set(c color.Color, v float64) color.Color {
	r, _, b, a := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(r), uint8(v), uint8(b), uint8(a)}
}

type blueCh struct{}

func (_ blueCh) Get(c color.Color) float64 {
	_, _, b, _ := utils.RatioRGBA(c)
	return b
}

func (_ blueCh) Set(c color.Color, v float64) color.Color {
	r, g, _, a := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(r), uint8(g), uint8(v), uint8(a)}
}

type alphaCh struct{}

func (_ alphaCh) Get(c color.Color) float64 {
	_, _, _, a := utils.RatioRGBA(c)
	return a
}

func (_ alphaCh) Set(c color.Color, v float64) color.Color {
	r, g, b, _ := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(v)}
}

type hueCh struct{}

func (_ hueCh) Get(c color.Color) float64 {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	return h.H / 360.0 // need value in range [0,1]
}

func (_ hueCh) Set(c color.Color, v float64) color.Color {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	h.H = math.Mod(v*360, 360) // again, need to scale from [0,1] to [0,360]
	return h
}

type saturationCh struct{}

func (_ saturationCh) Get(c color.Color) float64 {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	return h.S
}

func (_ saturationCh) Set(c color.Color, v float64) color.Color {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	h.S = v
	if h.S > 1 {
		h.S = 1
	} else if h.S < 0 {
		h.S = 0
	}
	return h
}

type lightnessCh struct{}

func (_ lightnessCh) Get(c color.Color) float64 {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	return h.L
}

func (_ lightnessCh) Set(c color.Color, v float64) color.Color {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	h.L = v
	if h.L > 1 {
		h.L = 1
	} else if h.L < 0 {
		h.L = 0
	}
	return h
}

type intensityCh struct{}

func (_ intensityCh) Get(c color.Color) float64 {
	h := altcolor.HSIAModel.Convert(c).(altcolor.HSIA)
	return h.I
}

func (_ intensityCh) Set(c color.Color, v float64) color.Color {
	h := altcolor.HSIAModel.Convert(c).(altcolor.HSIA)
	h.I = v
	if h.I > 1 {
		h.I = 1
	} else if h.I < 0 {
		h.I = 0
	}
	return h
}
