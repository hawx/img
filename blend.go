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

type UIntMixer (func(i, j uint32) uint32)
type FloatMixer (func (i, j float64) float64)

// Applies +f+ to each colour channel. Passes channel value as uint32, between 0
// and 255, it should return a value in the same range.
func EachChannel(c, d color.Color, fs... UIntMixer) color.Color {
	var f, q UIntMixer
	switch len(fs) {
	case 1:  f = fs[0]; q = fs[0]
	case 2:  f = fs[0]; q = fs[1]
	default: os.Exit(1)
	}

	i, j, k, l := utils.NormalisedRGBA(c)
	m, n, o, p := utils.NormalisedRGBA(d)

	r := utils.Truncate(f(i, m))
	g := utils.Truncate(f(j, n))
	b := utils.Truncate(f(k, o))
	a := utils.Truncate(q(l, p))

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Same as EachChannel but passes values as floats, and expects a float in
// return.
func EachChannelf(c, d color.Color, fs... FloatMixer) color.Color {
	var f, q FloatMixer
	switch len(fs) {
	case 1:  f = fs[0]; q = fs[0]
	case 2:  f = fs[0]; q = fs[1]
	default: os.Exit(1)
	}

	i, j, k, l := utils.NormalisedRGBAf(c)
	m, n, o, p := utils.NormalisedRGBAf(d)

	r := utils.Truncatef(f(i, m))
	g := utils.Truncatef(f(j, n))
	b := utils.Truncatef(f(k, o))
	a := utils.Truncatef(q(l, p))

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Applies +f+ to each colour channel. Passes channel values as float64, between
// 0 and 1, expects returned value to be within the same range.
func EachRatioChannel(c, d color.Color, fs... FloatMixer) color.Color {
	var f, q FloatMixer
	switch len(fs) {
	case 1:  f = fs[0]; q = fs[0]
	case 2:  f = fs[0]; q = fs[1]
	default: os.Exit(1)
	}

	i, j, k, l := utils.RatioRGBA(c)
	m, n, o, p := utils.RatioRGBA(d)

	r := utils.Truncatef(f(i, m) * 255)
	g := utils.Truncatef(f(j, n) * 255)
	b := utils.Truncatef(f(k, o) * 255)
	a := utils.Truncatef(q(l, p) * 255)

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
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



func fade(img image.Image, amount float64) image.Image {
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

func blend(a, b image.Image, opacity float64) image.Image {
	f := func(c, d color.Color) color.Color {
		return EachRatioChannel(c, d, func(i, j float64) float64 {
			return (i * (1 - opacity) + j * opacity)
		})
	}

	return BlendPixels(a, b, f)
}

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
			return i * j
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
				return (1 - (1 - i) * (1 - (j - 0.5)))
			}
			return (i * (j + 0.5))
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
