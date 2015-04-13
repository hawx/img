// Package shuffle implements functions to randomly swap pixels in an image.
package shuffle

import (
	"image"
	"math/rand"
)

func randBetween(a, b int) int {
	return rand.Intn(b - a)
}

// Shuffle moves the pixels in the Image to new, randomly chosen positions.
func Shuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			xb := randBetween(b.Min.X, b.Max.X)
			yb := randBetween(b.Min.Y, b.Max.Y)

			a := img.At(x, y)
			b := img.At(xb, yb)

			o.Set(x, y, b)
			o.Set(xb, yb, a)
		}
	}

	return o
}

// Vertically shuffles the pixels in the Image along the vertical axis only.
func Vertically(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			yb := randBetween(b.Min.Y, b.Max.Y)

			a := img.At(x, y)
			b := img.At(x, yb)

			o.Set(x, y, b)
			o.Set(x, yb, a)
		}
	}

	return o
}

// Horizontally shuffles the pixels in the Image along the horizontal axis only.
func Horizontally(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			xb := randBetween(b.Min.X, b.Max.X)

			a := img.At(x, y)
			b := img.At(xb, y)

			o.Set(x, y, b)
			o.Set(xb, y, a)
		}
	}

	return o
}
