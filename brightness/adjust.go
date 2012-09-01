package brightness

import (
	"../utils"
	"image"
	"image/color"
)

// Adjusts the brightness of +img+ by the given +value+. A +value+ of 0 has no
// effect.
func Adjust(img image.Image, value float64) image.Image {
	value  = (100 + value) / 100
	value *= value

	f := func(c color.Color) color.Color {
		r, g, b, a := utils.RatioRGBA(c)

		r = utils.Truncatef(r * value * 255)
		g = utils.Truncatef(g * value * 255)
		b = utils.Truncatef(b * value * 255)

		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a * 255)}
	}

	return utils.EachPixel(img, f)
}
