package utils

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"
)

func splitRectangle(b image.Rectangle, parts int) []image.Rectangle {
	if b.Dx() > b.Dy() {
		return ChopRectangle(b, parts, 1, ADD)
	}
	return ChopRectangle(b, 1, parts, ADD)
}

// EachColor iterates through each pixel of the Image, applying the function
// to each colour.
func EachColor(img image.Image, f func(c color.Color)) {
	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			f(img.At(x, y))
		}
	}
}

// PEachColor is like EachColor, but runs in parallel. This means that order can
// not be guaranteed.
func PEachColor(img image.Image, f func(c color.Color)) {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	c := make(chan int, nCPU)

	for _, r := range splitRectangle(img.Bounds(), nCPU) {
		go peachColorWorker(img, r, f, c)
	}

	// wait until work is done
	for i := 0; i < nCPU; i++ {
		<-c
	}
}

func peachColorWorker(img image.Image, b image.Rectangle, f func(c color.Color), c chan int) {

	EachColorInRectangle(img, b, f)
	c <- 1
}

// EachColorInRectangle is a helper function for working on a part of an image.
func EachColorInRectangle(img image.Image, b image.Rectangle, f func(c color.Color)) {
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			f(img.At(x, y))
		}
	}
}

// MapColor iterates through each pixel of the Image and applies the given
// function, drawing the returned colour to a new Image which is then returned.
func MapColor(img image.Image, f Composable) image.Image {
	// Use maximum number of CPUs available
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	c := make(chan int, nCPU)
	o := image.NewRGBA(img.Bounds())

	for _, r := range splitRectangle(img.Bounds(), nCPU) {
		go mapColorWorker(img, r, o, f, c)
	}

	// wait until work is done
	for i := 0; i < nCPU; i++ {
		<-c
	}

	return o
}

func mapColorWorker(img image.Image, bounds image.Rectangle, dest draw.Image,
	f Composable, c chan int) {

	MapColorInRectangle(img, bounds, dest, f)
	c <- 1
}

// MapColorInRectangle is a helper function for working on part of an image. It
// takes the original image, a function to use, a image to write to, and the
// bounds of the original (and therefore the final image) to act upon.
func MapColorInRectangle(img image.Image, bounds image.Rectangle, dest draw.Image,
	f Composable) {

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dest.Set(x, y, f(img.At(x, y)))
		}
	}
}
