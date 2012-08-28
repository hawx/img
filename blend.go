package main

import (
	"./utils"
	"os"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"image"
	"image/png"
	"image/color"
)

// http://docs.gimp.org/en/gimp-concepts-layer-modes.html

func Lightness(c color.Color) uint8 {
	r, g, b, _ := utils.NormalisedRGBA(c)

	maxi := utils.Max(r, g, b)
	mini := utils.Min(r, g, b)

	return uint8((maxi + mini) / 2)
}

func EachChannel(c, d color.Color, f (func(i, j uint32) uint32)) color.Color {
	i, j, k, l := utils.NormalisedRGBA(c)
	m, n, o, p := utils.NormalisedRGBA(d)

	r := utils.Truncate(f(i, m))
	g := utils.Truncate(f(j, n))
	b := utils.Truncate(f(k, o))
	a := utils.Truncate(f(l, p))

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func EachChannelf(c, d color.Color, f (func(i, j float64) float64)) color.Color {
	i, j, k, l := utils.NormalisedRGBAf(c)
	m, n, o, p := utils.NormalisedRGBAf(d)

	r := utils.Truncatef(f(i, m))
	g := utils.Truncatef(f(j, n))
	b := utils.Truncatef(f(k, o))
	a := utils.Truncatef(f(l, p))

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func EachRatioChannel(c, d color.Color, f (func(i, j float64) float64)) color.Color {
	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := utils.Truncatef(f(i, m))
	g := utils.Truncatef(f(j, n))
	b := utils.Truncatef(f(k, o))
	a := utils.Truncatef(f(l, p))

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

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


func blend(a, b image.Image, opacity float64) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachRatioChannel(c, d, func(i, j float64) float64 {
			return (i * (1 - opacity) + j * opacity) * 255
		})
	}

	return BlendPixels(a, b, f)
}

// +ib+ the base image
// +is+ the blend image
// +ir+ the result image from the blend operation
/*
 func properBlend(ib, is, ir image.Image) image.Image {
	bounds := ir.Bounds()
	result := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			cb := ib.At(x, y); cs := is.At(x, y); cr := ir.At(x, y)

			rb, gb, bb, ab := utils.RatioRGBA(cb)
			rs, gs, bs, as := utils.RatioRGBA(cs)
			rr, gr, br, ar := utils.RatioRGBA(cr)

			r := (1 - as/ar) * rb + as/ar * ((1 - ab) * rs + ab * rr)
			g := (1 - as/ar) * gb + as/ar * ((1 - ab) * gs + ab * gr)
			b := (1 - as/ar) * bb + as/ar * ((1 - ab) * bs + ab * br)
			a := (1 - as/ar) * ab + as/ar * ((1 - ab) * as + ab * ar)

		}
	}
}
 */


// Selects the blend colour for each pixel
func normal(a, b image.Image) image.Image {
	return b
}

// Randomly shows opacity percent of b's pixels.
func dissolve(a, b image.Image, opacity float64) image.Image {
	f := func(c, d color.Color) color.Color {
		if rand.Float64() > opacity {
			return c
		}
		return d
	}

	return BlendPixels(a, b, f)
}

// Selects the darker of each pixels' colour channels
func darken(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannel(c, d, func(i, j uint32) uint32 {
			return utils.Min(i, j)
		})
	}

	return BlendPixels(a, b, f)
}

// Multiples the base and blend image values
func multiply(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachRatioChannel(c, d, func(i, j float64) float64 {
			return (i * 255) * j
		})
	}

	return BlendPixels(a, b, f)
}

// Darkens the base colour to reflect the blend colour
func burn(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannelf(c, d, func(i, j float64) float64 {
			return 255 - (256 * (255 - i)) / (j + 1)
		})
	}

	return BlendPixels(a, b, f)
}

// Selects the lighter of each pixels' colour channels
func lighten(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannel(c, d, func(i, j uint32) uint32 {
			return utils.Max(i, j)
		})
	}

	return BlendPixels(a, b, f)
}

// Multiplies the complements of the base and blend colour values, then
// complements the result
func screen(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannelf(c, d, func(i, j float64) float64 {
			return 255 - ((255 - i) * (255 - j)) / 255
		})
	}

	return BlendPixels(a, b, f)
}

// Brightens the base colour to reflect the blend colour
func dodge(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannelf(c, d, func(i, j float64) float64 {
			return (256 * i) / ((255 - j) + 1)
		})
	}

	return BlendPixels(a, b, f)
}

// Multiplies or screens the colours, depeneding on the base colour.
func overlay(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannelf(c, d, func(i, j float64) float64 {
			return (i / 255) * (i + ((2 * j) / 255) * (255 - i))
		})
	}

	return BlendPixels(a, b, f)
}

func softLight(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachRatioChannel(c, d, func(i, j float64) float64 {
			if j > 0.5 {
				return (1 - (1 - i) * (1 - (j - 0.5))) * 255
			}
			return (i * (j + 0.5)) * 255
		})
	}

	return BlendPixels(a, b, f)
}

func hardLight(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannelf(c, d, func(i, j float64) float64 {
			if j > 128 {
				return 255 - ((255 - 2 * (j - 128)) * (255 - i)) / 256
			}
			return (2 * j * i) / 256
		})
	}

	return BlendPixels(a, b, f)
}


func difference(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		i, j, k, a := utils.NormalisedRGBAf(c)
		m, n, o, _ := utils.NormalisedRGBAf(d)

		r := math.Abs(i - m)
		g := math.Abs(j - n)
		b := math.Abs(k - o)

		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

	return BlendPixels(a, b, f)
}

func addition(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachChannel(c, d, func(i, j uint32) uint32 {
			return utils.Min(i + j, 255)
		})
	}

	return BlendPixels(a, b, f)
}

func subtraction(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		i, j, k, l := utils.NormalisedRGBA(c)
		m, n, o, _ := utils.NormalisedRGBA(d)

		r := utils.Truncate(i - m)
		g := utils.Truncate(j - n)
		b := utils.Truncate(k - o)

		if m > i { r = 0 }
		if n > j { g = 0 }
		if o > k { b = 0 }

		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(l)}
	}

	return BlendPixels(a, b, f)
}


func hue(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		_, s, l, a := utils.ToHSLA(c)
		h, _, _, _ := utils.ToHSLA(d)

		return utils.ToRGBA(h, s, l, a)
	}

	return BlendPixels(a, b, f)
}

func saturation(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		h, _, l, a := utils.ToHSLA(c)
		_, s, _, _ := utils.ToHSLA(d)

		return utils.ToRGBA(h, s, l, a)
	}

	return BlendPixels(a, b, f)
}

func colour(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		_, _, l, a := utils.ToHSLA(c)
		h, s, _, _ := utils.ToHSLA(d)

		return utils.ToRGBA(h, s, l, a)
	}

	return BlendPixels(a, b, f)
}

func luminosity(a, b image.Image) image.Image {
	f := func(c, d color.Color) color.Color {
		h, s, _, a := utils.ToHSLA(c)
		_, _, l, _ := utils.ToHSLA(d)

		return utils.ToRGBA(h, s, l, a)
	}

	return BlendPixels(a, b, f)
}



func printHelp() {
	msg := "Usage: compose <other> [opts]\n" +
		"\n" +
		"  Takes a png file from STDIN and blends it with the png file at <other>\n" +
		"  using the method chosen. The result is printed to STDOUT.\n" +
		"\n" +
		"  Note: 'blend colour' refers to the colour taken from <other>, while 'base \n" +
		"  color' refers to the colour taken from STDIN.\n" +
		"\n" +
		"  --opacity      # Opacity of blended image layer (default: 0.5)\n" +
		"\n" +
		"  BASIC\n" +
		"  --normal       # Paints pixels using <other> (default)\n" +
		"  --dissolve     # Paints pixels from <other> randomly, depending on opacity\n" +
		"\n" +
		"  DARKEN\n" +
		"  --darken       # Selects the darkest value for each colour channel\n" +
		"  --multiply     # Multiplies each colour channel\n" +
		"  --burn         # Darkens the base colour to increase contrast\n" +
		"\n" +
		"  LIGHTEN\n" +
		"  --lighten      # Selects the lightest value for each colour channel\n" +
		"  --screen       # Multiples the inverse of each colour channel\n" +
		"  --dodge        # Brightens the base colour to decrease contrast\n" +
		"\n" +
		"  CONTRAST\n" +
		"  --overlay      # Multiplies or screens the colours, depending on the base colour\n" +
		"  --soft-light   # Darkens or lightens the colours, depending on the blend colour\n" +
		"  --hard-light   # Multiplies or screens the colours, depending on the blend colour\n" +
		"\n" +
		"  COMPARATIVE\n" +
		"  --difference   # Finds the absolute difference between the base and blend colour\n" +
		"  --addition     # Adds the blend colour to the base colour\n" +
		"  --subtraction  # Subtracts the blend colour from the base colour\n" +
		"\n" +
		"  HSL\n" +
		"  --hue          # Uses just the hue of the blend colour\n" +
		"  --saturation   # Uses just the saturation of the blend colour\n" +
		"  --color        # Uses just the hue and saturation of the blend colour\n" +
		"  --luminosity   # Uses just the luminosity of the blend colour\n" +
		"\n" +
		"  --help         # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var (
	help         = flag.Bool("help", false, "")
	opacity      = flag.Float64("opacity", 0.5, "")

	// BASIC
	normalM      = flag.Bool("normal", false, "")
	dissolveM    = flag.Bool("dissolve", false, "")

	// DARKEN
	darkenM      = flag.Bool("darken", false, "")
	multiplyM    = flag.Bool("multiply", false, "")
	burnM        = flag.Bool("burn", false, "")

	// LIGHTEN
	lightenM     = flag.Bool("lighten", false, "")
	screenM      = flag.Bool("screen", false, "")
	dodgeM       = flag.Bool("dodge", false, "")

	// CONTRAST
	overlayM     = flag.Bool("overlay", false, "")
	softLightM   = flag.Bool("soft-light", false, "")
	hardLightM   = flag.Bool("hard-light", false, "")

	// COMPARATIVE
	differenceM  = flag.Bool("difference", false, "")
	additionM    = flag.Bool("addition", false, "")
	subtractionM = flag.Bool("subtraction", false, "")

	// HSL
	hueM         = flag.Bool("hue", false, "")
	saturationM  = flag.Bool("saturation", false, "")
	colorM       = flag.Bool("color", false, "")
	luminosityM  = flag.Bool("luminosity", false, "")
)

func main() {
	flag.Parse()
	if *help { printHelp() }

	a := utils.ReadStdin()

	path := flag.Args()[0]
	f, _ := os.Open(path)
	b, _ := png.Decode(f)
	var img image.Image

	if *dissolveM {
		img = dissolve(a, b, *opacity)

	} else if *darkenM {
		img = darken(a, b)
	} else if *multiplyM {
		img = multiply(a, b)
	} else if *burnM {
		img = burn(a, b)

	} else if *lightenM {
		img = lighten(a, b)
	} else if *screenM {
		img = screen(a, b)
	} else if *dodgeM {
		img = dodge(a, b)

	} else if *overlayM {
		img = overlay(a, b)
	} else if *softLightM {
		img = softLight(a, b)
	} else if *hardLightM {
		img = hardLight(a, b)

	} else if *differenceM {
		img = difference(a, b)
	} else if *additionM {
		img = addition(a, b)
	} else if *subtractionM {
		img = subtraction(a, b)

	} else if *hueM {
		img = hue(a, b)
	} else if *saturationM {
		img = saturation(a, b)
	} else if *colorM {
		img = colour(a, b)
	} else if *luminosityM {
		img = luminosity(a, b)

	} else {
		img = normal(a, b)
	}

	if !*dissolveM {
		img = blend(a, img, *opacity)
	}

	utils.WriteStdout(img)
}
