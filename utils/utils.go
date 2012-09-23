package utils

import (
	"fmt"
	"os"
	"image"
	"image/png"
	"image/color"
)

// ReadStdin reads a png file from standard input.
func ReadStdin() image.Image {
	if s,_ := os.Stdin.Stat(); s.Size() == 0 {
		fmt.Fprintln(os.Stderr, "Error: need to provide image by STDIN. See \"img help\".")
		os.Exit(1)
	}
	img, _ := png.Decode(os.Stdin)

	return img
}

// WriteStdout writes an Image to standard output.
func WriteStdout(img image.Image) {
	png.Encode(os.Stdout, img)
}

// Warn prints a message to standard error
func Warn(s interface{}) {
	fmt.Fprintln(os.Stderr, s)
}

// NormalisedRGBA returns the RGBA colour channel values of a Color in a
// normalised form.
func NormalisedRGBA(c color.Color) (rn, gn, bn, an uint32) {
	d := color.NRGBAModel.Convert(c).(color.NRGBA)
	r := d.R; g := d.G; b := d.B; a := d.A

	// Need to do some crazy type conversions first
	rn = uint32(uint8(r))
	gn = uint32(uint8(g))
	bn = uint32(uint8(b))
	an = uint32(uint8(a))

	return
}

// NormalisedRGBAf returns the RGBA colour channel values as floating point
// numbers with values from 0 to 255.
func NormalisedRGBAf(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := NormalisedRGBA(c)

	rn = float64(r)
	gn = float64(g)
	bn = float64(b)
	an = float64(a)

	return
}

// RatioRGBA returns the RGBA colour channel values as floating point numbers
// with values from 0 to 1.
func RatioRGBA(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := NormalisedRGBAf(c)

	rn = r / 255
	gn = g / 255
	bn = b / 255
	an = a / 255

	return
}

// Truncate takes a colour channel value and forces it into the range 0 to 255
// by setting any value below 0 to 0 and and any above 255 to 255.
func Truncate(n uint32) uint32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

// Truanctef is identical to Truncate but takes and returns a float64.
func Truncatef(n float64) float64 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

// Closeness calculates the "closeness" of two colours by finding the sum of
// differences in each colour channel.
func Closeness(one, two color.Color) uint32 {
	a, b, c, d := NormalisedRGBA(one)
	w, x, y, z := NormalisedRGBA(two)

	return (a - w) + (b - x) + (c - y) + (d - z)
}

// Average takes a list of colours and returns the average.
func Average(cs... color.Color) color.Color {
	var red, green, blue, alpha uint32
	red = 0; green = 0; blue = 0; alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := NormalisedRGBA(cs[i])

		red += r; green += g; blue += b; alpha += a
	}

	return color.NRGBA{
		uint8(red   / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue  / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}

// Min returns the smallest value in the list of uint32s given.
func Min(ns... uint32) (n uint32) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] < n {
			n = ns[i]
		}
	}
	return
}

// Minf returns the smallest value in the list of float64s given.
func Minf(ns... float64) (n float64) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] < n {
			n = ns[i]
		}
	}
	return
}

// Max returns the largest value in the list of uint32s given.
func Max(ns... uint32) (n uint32) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] > n {
			n = ns[i]
		}
	}
	return
}

// Maxf returns the largest value in the list of float64s given.
func Maxf(ns... float64) (n float64) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] > n {
			n = ns[i]
		}
	}
	return
}

// EachPixel iterates through each pixel of the Image applies the function given
// and draws the result to a new Image which is then returned.
func EachPixel(img image.Image, f func(c color.Color) color.Color) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			result := f(img.At(x, y))
			o.Set(x, y, result)
		}
	}

	return o
}
