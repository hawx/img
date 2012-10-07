package channel

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
)

// Red applies the Adjuster to the red channel of each pixel in the Image.
func Red(img image.Image, adj utils.Adjuster) image.Image {
	f := func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		r = adj(r)
		if r > 1 { r = 1 } else if r < 0 { r = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}

	return utils.EachPixel(img, f)
}

// Green applies the Adjuster to the green channel of each pixel in the Image.
func Green(img image.Image, adj utils.Adjuster) image.Image {
	f := func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		g = adj(g)
		if g > 1 { g = 1 } else if g < 0 { g = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}

	return utils.EachPixel(img, f)
}

// Blue applies the Adjuster to the blue channel of each pixel in the Image.
func Blue(img image.Image, adj utils.Adjuster) image.Image {
	f := func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		b = adj(b)
		if b > 1 { b = 1 } else if b < 0 { b = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}

	return utils.EachPixel(img, f)
}

// Alpha applies the Adjuster to the alpha channel of each pixel in the Image.
func Alpha(img image.Image, adj utils.Adjuster) image.Image {
	f := func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)
		a = adj(a)
		if a > 1 { a = 1 } else if a < 0 { a = 0 }
		return color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), uint8(a * 255)}
	}

	return utils.EachPixel(img, f)
}
