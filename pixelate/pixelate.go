// Package pixelate implements various functions for pixelating images.
package pixelate

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"

	"hawx.me/code/img/utils"
)

type Style int

const (
	CROPPED Style = iota
	FITTED
)

func paintAverage(img image.Image, bounds image.Rectangle, dest draw.Image,
	c chan int) {
	values := make([]color.Color, bounds.Dx()*bounds.Dy())
	count := 0

	utils.EachColorInRectangle(img, bounds, func(c color.Color) {
		values[count] = c
		count++
	})

	avg := utils.Average(values...)

	utils.MapColorInRectangle(img, bounds, dest, func(c color.Color) color.Color {
		return avg
	})

	c <- 1
}

// Pixelate takes an Image and pixelates it into rectangles with the dimensions
// given. The colour values in each region are averaged to produce the resulting
// colours.
func Pixelate(img image.Image, size utils.Dimension, style Style) image.Image {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	var o draw.Image
	b := img.Bounds()
	c := make(chan int, nCPU)
	i := 0 // like pxl. this may not be the best way

	switch style {
	case CROPPED:
		cols := b.Dx() / size.W
		rows := b.Dy() / size.H

		o = image.NewRGBA(image.Rect(0, 0, size.W*cols, size.H*rows))

		for j, r := range utils.ChopRectangleToSizes(b, size.H, size.W, utils.IGNORE) {
			go paintAverage(img, r, o, c)
			i = j
		}

	case FITTED:
		o = image.NewRGBA(b)

		for j, r := range utils.ChopRectangleToSizes(b, size.H, size.W, utils.SEPARATE) {
			go paintAverage(img, r, o, c)
			i = j
		}
	}

	for j := 0; j < i; j++ {
		<-c
	}

	return o
}
