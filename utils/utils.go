package utils

import (
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

func NormalisedRGBA(c color.Color) (rn, gn, bn, an uint32) {
	r, g, b, a := c.RGBA()

	// Need to do some crazy type conversions first
	rn = uint32(uint8(r))
	gn = uint32(uint8(g))
	bn = uint32(uint8(b))
	an = uint32(uint8(a))

	return
}

func RatioRGBA(c color.Color) (rn, gn, bn, an float32) {
	r, g, b, a := c.RGBA()

	rn = float32(uint8(r)) / 255
	gn = float32(uint8(g)) / 255
	bn = float32(uint8(b)) / 255
	an = float32(uint8(a)) / 255

	return
}

func TruncateInt(n uint32) uint32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

func TruncateFloat(n float32) float32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

func Closeness(one, two color.Color) uint32 {
	a, b, c, d := NormalisedRGBA(one)
	w, x, y, z := NormalisedRGBA(two)

	return (a - w) + (b - x) + (c - y) + (d - z)
}

func Average(cs []color.Color) color.Color {
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
