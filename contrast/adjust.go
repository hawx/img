package contrast

import (
	"../utils"
	"image"
	"image/color"
)

// Adjusts the contrast of +img+ by given +value+. A +value+ of 0 has no effect.
func Adjust(img image.Image, value float64) image.Image {
	value  = (100 + value) / 100
	value *= value

	f := func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)

		r = utils.Truncatef((((r - 0.5) * value) + 0.5) * 255)
		g = utils.Truncatef((((g - 0.5) * value) + 0.5) * 255)
		b = utils.Truncatef((((b - 0.5) * value) + 0.5) * 255)
		a = a * 255

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

	return utils.EachPixel(img, f)
}
