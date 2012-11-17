package levels

import (
	"github.com/hawx/img/utils"
	"image/color"
)

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
	RedChannel   Channel = createChannel(getRed, setRed)
	GreenChannel Channel = createChannel(getGreen, setGreen)
	BlueChannel  Channel = createChannel(getBlue, setBlue)
	AlphaChannel Channel = createChannel(getAlpha, setAlpha)
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
