package blend

import (
	"../utils"
	"math/rand"
	"image"
	"image/color"
)


func BlendPixels(a, b image.Image, f func(c, d color.Color) color.Color) image.Image {
	ba := a.Bounds(); bb := b.Bounds()
	width  := int(utils.Min(uint32(ba.Dx()), uint32(bb.Dx())))
	height := int(utils.Min(uint32(ba.Dy()), uint32(bb.Dy())))

	result := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := f(a.At(x, y), b.At(x, y))
			result.Set(x, y, pixel)
		}
	}

	return result
}


func Fade(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		r, g, b, a := utils.NormalisedRGBA(c)

		return color.NRGBA{
			uint8(float64(r)),
			uint8(float64(g)),
			uint8(float64(b)),
			uint8(float64(a) * amount),
		}
	}

	return utils.EachPixel(img, f)
}

func ratioNRGBA(r, g, b, a float64) color.Color {
	return color.NRGBA{
		uint8(utils.Truncatef(r * 255)),
		uint8(utils.Truncatef(g * 255)),
		uint8(utils.Truncatef(b * 255)),
		uint8(utils.Truncatef(a * 255)),
	}
}


func normal(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := i * (l - p) + m * p
	g := j * (l - p) + n * p
	b := k * (l - p) + o * p
	a := l + p

	return ratioNRGBA(r, g, b, a)
}

func dissolve(c, d color.Color) color.Color {
	var r, g, b, a float64
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	if rand.Float64() < p {
		r = m; g = n; b = o; a = 1
	} else {
		r = i; g = j; b = k; a = utils.Maxf(l, p)
	}

	return ratioNRGBA(r, g, b, a)
}

func Normal(a, b image.Image) image.Image {
	return BlendPixels(a, b, normal)
}

func Dissolve(a, b image.Image) image.Image {
	return BlendPixels(a, b, dissolve)
}
