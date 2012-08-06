package utils

import (
	"fmt"
	"os"
	"image"
	"image/png"
	"image/color"
)

func ReadStdin() image.Image {
	img, _ := png.Decode(os.Stdin)

	return img
}

func WriteStdout(img image.Image) {
	png.Encode(os.Stdout, img)
}

func Warn(s interface{}) {
	fmt.Fprintln(os.Stderr, s)
}

func NormalisedRGBA(c color.Color) (rn, gn, bn, an uint32) {
	r, g, b, a := c.RGBA()

	// Need to do some crazy type conversions first
	rn = uint32(uint8(r))
	gn = uint32(uint8(g))
	bn = uint32(uint8(b))
	an = uint32(uint8(a))

	return
}

func RatioRGBA(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := c.RGBA()

	rn = float64(uint8(r)) / 255
	gn = float64(uint8(g)) / 255
	bn = float64(uint8(b)) / 255
	an = float64(uint8(a)) / 255

	return
}

func Truncate(n uint32) uint32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

func Truncatef(n float64) float64 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

func Closeness(one, two color.Color) uint32 {
	a, b, c, d := NormalisedRGBA(one)
	w, x, y, z := NormalisedRGBA(two)

	return (a - w) + (b - x) + (c - y) + (d - z)
}

func Average(cs... color.Color) color.Color {
	var red, green, blue, alpha uint32
	red = 0; green = 0; blue = 0; alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := NormalisedRGBA(cs[i])

		red += r; green += g; blue += b; alpha += a
	}

	return color.RGBA{
		uint8(red   / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue  / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}

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

func ChangePixels(img image.Image, f func(c color.Color) color.Color) image.Image {
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
