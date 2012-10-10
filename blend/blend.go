// Package blend implements various blending mode functions between two
// images. These are modelled on those in Adobe Photoshop. Each function takes
// as the first argument the "base image" and the second the "blend image", they
// then return an image blending these in specific ways.
package blend

import (
	"github.com/hawx/img/utils"
	"github.com/hawx/img/hsla"
	"math"
	"math/rand"
	"image"
	"image/color"
)

// Blender takes two colours (base and blend, respectively) and returns another
// colour.
type Blender (func (c, d color.Color) color.Color)

// BlendPixels takes the base and blend images and applies the given Blender to
// each of their pixel pairs.
func BlendPixels(a, b image.Image, f Blender) image.Image {
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

// Fade changes the opacity of the Image given by the amount given. The
// resulting opacity is the product of the image's opacity and the amount, so a
// value of 1 has no effect whilst a value of 0 makes the image fully
// transparent.
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

	return utils.MapColor(img, f)
}

func ratioNRGBA(r, g, b, a float64) color.Color {
	return color.NRGBA{
		uint8(utils.Truncatef(r * 255)),
		uint8(utils.Truncatef(g * 255)),
		uint8(utils.Truncatef(b * 255)),
		uint8(utils.Truncatef(a * 255)),
	}
}


// Normal selects the blend Image.
func Normal(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		return d
	})
}

// Dissolve randomly selects pixels from the blend image, depending on their
// opacity. Blend pixels with higher opacities are more likely to be displayed.
func Dissolve(a, b image.Image) image.Image {
	ba := a.Bounds(); bb := b.Bounds()
	width  := int(utils.Min(uint32(ba.Dx()), uint32(bb.Dx())))
	height := int(utils.Min(uint32(ba.Dy()), uint32(bb.Dy())))

	result := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// base colour
			rb, gb, bb, ab := utils.RatioRGBA(a.At(x, y))
			// blend colour
			rs, gs, bs, as := utils.RatioRGBA(b.At(x, y))

			toPaint := ratioNRGBA(rb, gb, bb, ab)

			if rand.Float64() < as {
				toPaint = ratioNRGBA(rs, gs, bs, 1)
			}

			result.Set(x, y, toPaint)
		}
	}

	return result
}

// Darken selects the darkest value for each pixels' colour channels.
func Darken(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := utils.Minf(i, m)
		g := utils.Minf(j, n)
		b := utils.Minf(k, o)
		a := utils.Minf(l, p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Multiply multiplies the base and blend image colour channels.
func Multiply(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := i * m
		g := j * n
		b := k * o
		a := l * p

		return ratioNRGBA(r, g, b, a)
	})
}

// Burn darkens the base colour to reflect the blend colour.
func Burn(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := 1 - ((1 - i) / m)
		g := 1 - ((1 - j) / n)
		b := 1 - ((1 - k) / o)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// LinearBurn adds the values of each colour channel together, then subtracts
// white to produce a darker image.
func LinearBurn(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := i + m - 1
		g := j + n - 1
		b := k + o - 1
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Darker chooses the darkest colour by comparing the sum of the colour channels.
func Darker(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, _ := utils.RatioRGBA(c)
		m, n, o, _ := utils.RatioRGBA(d)

		if i + j + k < m + n + o {
			return c
		}
		return d
	})
}

// Lightne selects the lighter of each pixels' colour channels.
func Lighten(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := utils.Maxf(i, m)
		g := utils.Maxf(j, n)
		b := utils.Maxf(k, o)
		a := utils.Maxf(l, p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Screen multiplies the complements of the base and blend colour channel
// values, then complements the result.
func Screen(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := 1 - ((1 - i) * (1 - m))
		g := 1 - ((1 - j) * (1 - n))
		b := 1 - ((1 - k) * (1 - o))
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Dodge brightens the base colour to reflect the blend colour.
func Dodge(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := i / (1 - m)
		g := j / (1 - n)
		b := k / (1 - o)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// LinearDodge adds the values for each colour channel together.
func LinearDodge(a, b image.Image) image.Image {
	return Addition(a, b)
}

// Lighter chooses the lightest colour by comparing the sum of the colour
// channels.
func Lighter(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, _ := utils.RatioRGBA(c)
		m, n, o, _ := utils.RatioRGBA(d)

		if i + j + k > m + n + o {
			return c
		}
		return d
	})
}

// Overlay multiplies or screens the colours, depending on the base colour.
func Overlay(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
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
	})
}

// SoftLight darkens or lightens the colours, depending on the blend colour. The
// effect is similar to shining a soft spotlight on the image.
func SoftLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
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
	})
}

// HardLight multiplies or screens the colours, depending on the blend
// colour. The effect is similar to shining a harsh spotlight on the image.
func HardLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
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
	})
}

// VividLight combines Dodge and Burn. Dodge applies to lighter colours, and
// Burn to darker.
func VividLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		f := func(i, j float64) float64 {
			if j > 0.5 {
				return i / (2 * (1 - j))
			}
			return 1 - (1 - i) / (2 * j)
		}

		r := f(i, m)
		g := f(j, n)
		b := f(k, o)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// LinearLight lightens or darkens the image by changing the brightness. If the
// blend colour is lighter, the image is lightened; if the blend colour is
// darker, the image is darkened. It uses linear burn and linear dodge to darken
// or lighten.
func LinearLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		f := func(i, j float64) float64 {
			if j > 0.5 {
				return i + 2 * (j - 0.5)
			}
			return i + 2 * j - 1
		}

		r := f(i, m)
		g := f(j, n)
		b := f(k, o)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// PinLight replaces the colours, depending on the blend colour.
func PinLight(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		f := func(i, j float64) float64 {
			if i < 2 * j - 1 {
				return 2 * j - 1
			} else if i > 2 * j  {
				return 2 * j
			}
			return i
		}

		r := f(i, m)
		g := f(j, n)
		b := f(k, o)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// HardMix adds the red, green and blue channel values of the blend colour to
// the RGB values of the base colour. It sets any values greater than 255 to
// 255, and anything less to 0. This therefore makes all pixels either red,
// green, blue, white or black.
func HardMix(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		f := func(i, j float64) float64 {
			if j < 1 - i {
				return 0
			}
			return 1
		}

		r := f(i, m)
		g := f(j, n)
		b := f(k, o)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Difference finds the absolute difference between the base and blend colours.
func Difference(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := math.Abs(m - i)
		g := math.Abs(n - j)
		b := math.Abs(o - k)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Exclusion creates an effect similar to, but lower in contrast than,
// difference.
func Exclusion(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.RatioRGBA(c)
		m, n, o, p := utils.RatioRGBA(d)

		r := m + i - (2 * m * i)
		g := n + j - (2 * n * j)
		b := o + k - (2 * o * k)
		a := p + l * (1 - p)

		return ratioNRGBA(r, g, b, a)
	})
}

// Addition adds the blend colour to the base colour. (aka. Linear Dodge)
func Addition(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i, j, k, l := utils.NormalisedRGBA(c)
		m, n, o, p := utils.NormalisedRGBA(d)

		r := utils.Min(i + m, 255)
		g := utils.Min(j + n, 255)
		b := utils.Min(k + o, 255)
		a := utils.Min(l + p, 255)

		return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	})
}

// Subtraction subtracts the blend colour from the base colour.
func Subtraction(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
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
	})
}

// Hue uses the hue of the blend colour, with the saturation and luminosity of
// the base colour.
func Hue(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i := hsla.HSLAModel.Convert(c).(hsla.HSLA)
		j := hsla.HSLAModel.Convert(d).(hsla.HSLA)
		i.H = j.H

		return i
	})
}

// Saturation uses the saturation of the blend colour, with the hue and
// luminosity of the base colour.
func Saturation(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i := hsla.HSLAModel.Convert(c).(hsla.HSLA)
		j := hsla.HSLAModel.Convert(d).(hsla.HSLA)
		i.S = j.S

		return i
	})
}

// Color uses the hue and saturation of the blend colour, with the luminosity of
// the base colour.
func Color(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i := hsla.HSLAModel.Convert(c).(hsla.HSLA)
		j := hsla.HSLAModel.Convert(d).(hsla.HSLA)
		i.H = j.H
		i.S = j.S

		return i
	})
}

// Luminosity uses the luminosity of the blend colour, with the hue and
// saturation of the base colour.
func Luminosity(a, b image.Image) image.Image {
	return BlendPixels(a, b, func (c, d color.Color) color.Color {
		i := hsla.HSLAModel.Convert(c).(hsla.HSLA)
		j := hsla.HSLAModel.Convert(d).(hsla.HSLA)
		i.L = j.L

		return i
	})
}
