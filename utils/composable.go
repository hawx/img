package utils

import (
	"image"
	"image/color"
)

// A Composable function is one which can be composed with another function of
// the same type. It simply transforms a Color into another Color.
type Composable func(color.Color) color.Color

// Map takes a Composable function, and returns a function that applies it to
// every pixel of an Image.
func Map(f func() Composable) func(image.Image) image.Image {
	return func(img image.Image) image.Image {
		return MapColor(img, f())
	}
}

// MapAdjuster takes a Composable function and an Adjuster, and returns a
// function that applies the function to every pixel of an Image.
func MapAdjuster(f func(Adjuster) Composable) func(image.Image, Adjuster) image.Image {
	return func(img image.Image, adj Adjuster) image.Image {
		return MapColor(img, f(adj))
	}
}

// Compose takes a variable list of Composable functions and returns a single
// Composable function which performs them sequentially. For example,
//
//   var f Composable = Compose(
//     greyscale.BlueC(),
//     brightness.AdjustC(utils.Adder(0.05)),
//   )
//
//   // Only loops through image once!
//   img = MapColors(img, f)
//
func Compose(fs ...Composable) Composable {
	return func(c color.Color) color.Color {
		for _, f := range fs {
			c = f(c)
		}
		return c
	}
}
