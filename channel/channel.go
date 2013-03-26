// Package channel provides functions to alter the values of separate color
// channels in an image.
package channel

import (
	"github.com/hawx/img/utils"
	"github.com/hawx/img/altcolor"
	"image"
	"image/color"
	"math"
)

// Adjust applies the given Adjuster to the Image on only the Channels specified.
func Adjust(img image.Image, adj utils.Adjuster, chs... Channel) image.Image {
	return utils.MapColor(img, AdjustC(adj, chs...))
}

// AdjustC returns a Composable function that applies the given Adjuster to the
// Channels.
func AdjustC(adj utils.Adjuster, chs... Channel) utils.Composable {
	return func(c color.Color) color.Color {
		for _, ch := range chs {
			v := ch.Get(c)
			c  = ch.Set(c, adj(v))
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

type channelFuncs struct {
	g func(color.Color) float64
	s func(color.Color, float64) color.Color
}

func (c *channelFuncs) Get(o color.Color) float64 {
	return c.g(o)
}

func (c *channelFuncs) Set(o color.Color, v float64) color.Color {
	return c.s(o, v)
}

func createChannel(g func(color.Color) float64,
	s func(color.Color, float64) color.Color) Channel {

	return &channelFuncs{g, s}
}

var (
	Red        Channel = createChannel(getRed, setRed)
	Green      Channel = createChannel(getGreen, setGreen)
	Blue       Channel = createChannel(getBlue, setBlue)
	Alpha      Channel = createChannel(getAlpha, setAlpha)
	Hue        Channel = createChannel(getHue, setHue)
	Saturation Channel = createChannel(getSaturation, setSaturation)
	Lightness  Channel = createChannel(getLightness, setLightness)
	Intensity  Channel = createChannel(getIntensity, setIntensity)

	// Alias
	Brightness = Intensity
)

func getRed(c color.Color) float64 {
	r,_,_,_ := utils.RatioRGBA(c)
	return r
}

func setRed(c color.Color, v float64) color.Color {
	_,g,b,a := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(v), uint8(g), uint8(b), uint8(a)}
}

func getGreen(c color.Color) float64 {
	_,g,_,_ := utils.RatioRGBA(c)
	return g
}

func setGreen(c color.Color, v float64) color.Color {
	r,_,b,a := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(r), uint8(v), uint8(b), uint8(a)}
}

func getBlue(c color.Color) float64 {
	_,_,b,_ := utils.RatioRGBA(c)
	return b
}

func setBlue(c color.Color, v float64) color.Color {
	r,g,_,a := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(r), uint8(g), uint8(v), uint8(a)}
}

func getAlpha(c color.Color) float64 {
	_,_,_,a := utils.RatioRGBA(c)
	return a
}

func setAlpha(c color.Color, v float64) color.Color {
	r,g,b,_ := utils.NormalisedRGBA(c)
	v = utils.Truncatef(255 * v)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(v)}
}

func getHue(c color.Color) float64 {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	return h.H / 360.0 // need value in range [0,1]
}

func setHue(c color.Color, v float64) color.Color {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	h.H = math.Mod(v * 360, 360) // again, need to scale from [0,1] to [0,360]
	return h
}

func getSaturation(c color.Color) float64 {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	return h.S
}

func setSaturation(c color.Color, v float64) color.Color {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	h.S = v
	if h.S > 1 { h.S = 1 } else if h.S < 0 { h.S = 0 }
	return h
}

func getLightness(c color.Color) float64 {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	return h.L
}

func setLightness(c color.Color, v float64) color.Color {
	h := altcolor.HSLAModel.Convert(c).(altcolor.HSLA)
	h.L = v
	if h.L > 1 { h.L = 1 } else if h.L < 0 { h.L = 0 }
	return h
}

func getIntensity(c color.Color) float64 {
	h := altcolor.HSIAModel.Convert(c).(altcolor.HSIA)
	return h.I
}

func setIntensity(c color.Color, v float64) color.Color {
	h := altcolor.HSIAModel.Convert(c).(altcolor.HSIA)
	h.I = v
	if h.I > 1 { h.I = 1 } else if h.I < 0 { h.I = 0 }
	return h
}
