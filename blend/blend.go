package blend

import (
	"../utils"
	"../hsl"
	"math"
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
			// Uses methods described in "PDF Reference, Third Edition" from Adobe
			//  see: http://www.adobe.com/devnet/pdf/pdf_reference_archive.html

			// backdrop colour
			cb := a.At(x, y)
			// source colour
			cs := b.At(x, y)
			// result colour
			cr := f(cb, cs)

			rb, gb, bb, ab := utils.RatioRGBA(cb)
			rs, gs, bs, as := utils.RatioRGBA(cs)
			rr, gr, br, _  := utils.RatioRGBA(cr)

			// Color compositing formula, expanded form. (Section 7.2.5)
			red   := ((1 - as) * ab * rb) + ((1 - ab) * as * rs) + (ab * as * rr)
			green := ((1 - as) * ab * gb) + ((1 - ab) * as * gs) + (ab * as * gr)
			blue  := ((1 - as) * ab * bb) + ((1 - ab) * as * bs) + (ab * as * br)

			// Union function. (Section 7.2.6)
			alpha := ab + as - (ab * as)

			result.Set(x, y, color.RGBA{
				uint8(utils.Truncatef(red * 255)),
				uint8(utils.Truncatef(green * 255)),
				uint8(utils.Truncatef(blue * 255)),
				uint8(utils.Truncatef(alpha * 255)),
			})
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
	return d
}

func dissolve(c, d color.Color) color.Color {
	if r,g,b,a := utils.RatioRGBA(d); rand.Float64() < a {
		return ratioNRGBA(r, g, b, 1)
	}
	return color.Transparent
}

func darken(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := utils.Minf(i, m)
	g := utils.Minf(j, n)
	b := utils.Minf(k, o)
	a := utils.Minf(l, p)

	return ratioNRGBA(r, g, b, a)
}

func multiply(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := i * m
	g := j * n
	b := k * o
	a := l * p

	return ratioNRGBA(r, g, b, a)
}

func burn(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := 1 - ((1 - i) / m)
	g := 1 - ((1 - j) / n)
	b := 1 - ((1 - k) / o)
	a := p + l * (1 - p)

	return ratioNRGBA(r, g, b, a)
}

func darker(c, d color.Color) color.Color {
	i, j, k, _ := utils.RatioRGBA(c)
	m, n, o, _ := utils.RatioRGBA(d)

	if i + j + k < m + n + o {
		return c
	}
	return d
}

func lighten(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := utils.Maxf(i, m)
	g := utils.Maxf(j, n)
	b := utils.Maxf(k, o)
	a := utils.Maxf(l, p)

	return ratioNRGBA(r, g, b, a)
}

func screen(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := 1 - ((1 - i) * (1 - m))
	g := 1 - ((1 - j) * (1 - n))
	b := 1 - ((1 - k) * (1 - o))
	a := p + l * (1 - p)

	return ratioNRGBA(r, g, b, a)
}

func dodge(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := i / (1 - m)
	g := j / (1 - n)
	b := k / (1 - o)
	a := p + l * (1 - p)

	return ratioNRGBA(r, g, b, a)
}

func lighter(c, d color.Color) color.Color {
	i, j, k, _ := utils.RatioRGBA(c)
	m, n, o, _ := utils.RatioRGBA(d)

	if i + j + k > m + n + o {
		return c
	}
	return d
}

func overlay(c, d color.Color) color.Color {
	i, j, k, l := utils.NormalisedRGBAf(c)
	m, n, o, p := utils.NormalisedRGBAf(d)

	r := (i / 255) * (i + ((2 * m) / 255) * (255 - i))
	g := (j / 255) * (j + ((2 * n) / 255) * (255 - j))
	b := (k / 255) * (k + ((2 * o) / 255) * (255 - k))
	a := p + l * (1 - p)

	return color.NRGBA{
		uint8(utils.Truncatef(r)),
		uint8(utils.Truncatef(g)),
		uint8(utils.Truncatef(b)),
		uint8(utils.Truncatef(a * 255)),
	}
}

func softLight(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	f := func(i, j float64) float64 {
		if j > 0.5 {
			return 1 - (1 - i) * (1 - (j - 0.5))
		}
		return i * (j + 0.5)
	}

	r := f(i, m)
	g := f(j, n)
	b := f(k, o)
	a := p + l * (1 - p)

	return ratioNRGBA(r, g, b, a)
}

func hardLight(c, d color.Color) color.Color {
	i, j, k, l := utils.NormalisedRGBAf(c)
	m, n, o, p := utils.NormalisedRGBAf(d)

	f := func(i, j float64) float64 {
		if j > 128 {
			return 255 - ((255 - 2 * (j - 128)) * (255 - i)) / 256
		}
		return (2 * j * i) / 256
	}

	r := f(i, m)
	g := f(j, n)
	b := f(k, o)
	a := p + l * (1 - p)

	return color.NRGBA{
		uint8(utils.Truncatef(r)),
		uint8(utils.Truncatef(g)),
		uint8(utils.Truncatef(b)),
		uint8(utils.Truncatef(a * 255)),
	}
}

func difference(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := math.Abs(m - i)
	g := math.Abs(n - j)
	b := math.Abs(o - k)
	a := p + l * (1 - p)

	return ratioNRGBA(r, g, b, a)
}

func exclusion(c, d color.Color) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := m + i - (2 * m * i)
	g := n + j - (2 * n * j)
	b := o + k - (2 * o * k)
	a := p + l * (1 - p)

	return ratioNRGBA(r, g, b, a)
}

// aka Linear Dodge
func addition(c, d color.Color) color.Color {
	i, j, k, l := utils.NormalisedRGBA(c)
	m, n, o, p := utils.NormalisedRGBA(d)

	r := utils.Min(i + m, 255)
	g := utils.Min(j + n, 255)
	b := utils.Min(k + o, 255)
	a := utils.Min(l + p, 255)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func subtraction(c, d color.Color) color.Color {
	i, j, k, l := utils.NormalisedRGBA(c)
	m, n, o, p := utils.NormalisedRGBA(d)

	r := utils.Truncate(i - m)
	g := utils.Truncate(j - n)
	b := utils.Truncate(k - o)

	if m > i { r = 0 }
	if n > j { g = 0 }
	if o > k { b = 0 }

	a := p + l * (1 - p)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func hue(c, d color.Color) color.Color {
	i := hsl.HSLAModel.Convert(c).(hsl.HSLA)
	j := hsl.HSLAModel.Convert(d).(hsl.HSLA)
	i.H = j.H

	return i
}

func saturation(c, d color.Color) color.Color {
	i := hsl.HSLAModel.Convert(c).(hsl.HSLA)
	j := hsl.HSLAModel.Convert(d).(hsl.HSLA)
	i.S = j.S

	return i
}

func colour(c, d color.Color) color.Color {
	i := hsl.HSLAModel.Convert(c).(hsl.HSLA)
	j := hsl.HSLAModel.Convert(d).(hsl.HSLA)
	i.H = j.H
	i.S = j.S

	return i
}

func luminosity(c, d color.Color) color.Color {
	i := hsl.HSLAModel.Convert(c).(hsl.HSLA)
	j := hsl.HSLAModel.Convert(d).(hsl.HSLA)
	i.L = j.L

	return i
}

func Normal(a, b image.Image) image.Image {
	return BlendPixels(a, b, normal)
}

func Dissolve(a, b image.Image) image.Image {
	return BlendPixels(a, b, dissolve)
}

func Darken(a, b image.Image) image.Image {
	return BlendPixels(a, b, darken)
}

func Multiply(a, b image.Image) image.Image {
	return BlendPixels(a, b, multiply)
}

func Burn(a, b image.Image) image.Image {
	return BlendPixels(a, b, burn)
}

func Darker(a, b image.Image) image.Image {
	return BlendPixels(a, b, darker)
}

func Lighten(a, b image.Image) image.Image {
	return BlendPixels(a, b, lighten)
}

func Screen(a, b image.Image) image.Image {
	return BlendPixels(a, b, screen)
}

func Dodge(a, b image.Image) image.Image {
	return BlendPixels(a, b, dodge)
}

func Lighter(a, b image.Image) image.Image {
	return BlendPixels(a, b, lighter)
}

func Overlay(a, b image.Image) image.Image {
	return BlendPixels(a, b, overlay)
}

func SoftLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, softLight)
}

func HardLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, hardLight)
}

func Difference(a, b image.Image) image.Image {
	return BlendPixels(a, b, difference)
}

func Exclusion(a, b image.Image) image.Image {
	return BlendPixels(a, b, exclusion)
}

func Addition(a, b image.Image) image.Image {
	return BlendPixels(a, b, addition)
}

func Subtraction(a, b image.Image) image.Image {
	return BlendPixels(a, b, subtraction)
}

func Hue(a, b image.Image) image.Image {
	return BlendPixels(a, b, hue)
}

func Saturation(a, b image.Image) image.Image {
	return BlendPixels(a, b, saturation)
}

func Color(a, b image.Image) image.Image {
	return BlendPixels(a, b, colour)
}

func Luminosity(a, b image.Image) image.Image {
	return BlendPixels(a, b, luminosity)
}
