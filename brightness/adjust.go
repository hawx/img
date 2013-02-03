// Package brightness implements a single function to adjust the brightness of
// an image.
package brightness

import (
	"github.com/hawx/img/utils"
	"image/color"
)

var Adjust = utils.MapAdjuster(AdjustC)

func AdjustC(adj utils.Adjuster) utils.Composable {
	return func(c color.Color) color.Color {
		r, g, b, a := utils.RatioRGBA(c)

		r = utils.Truncatef(adj(r) * 255)
		g = utils.Truncatef(adj(g) * 255)
		b = utils.Truncatef(adj(b) * 255)

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a * 255)}
	}
}
