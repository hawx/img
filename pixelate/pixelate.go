// Package pixelate implements various functions for pixelating images.
package pixelate

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
)

// Pixelate takes an Image and pixelates it into rectangles with the dimensions
// given. The colour values in each region are averaged to produce the resulting
// colours.
func Pixelate(img image.Image, size utils.Dimension) image.Image {
	b := img.Bounds()

	cols := b.Dx() / size.W
	rows := b.Dy() / size.H

	o := image.NewRGBA(image.Rect(0, 0, size.W * cols, size.H * rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			values := make([]color.Color, size.H * size.W)
			count  := 0

			for y := 0; y < size.H; y++ {
				for x := 0; x < size.W; x++ {

					realY := row * size.H + y
					realX := col * size.W + x

					values[count] = img.At(realX, realY)
					count++
				}
			}

			avg := utils.Average(values...)

			for y := 0; y < size.H; y++ {
				for x := 0; x < size.W; x++ {

					realY := row * size.H + y
					realX := col * size.W + x

					o.Set(realX, realY, avg)
				}
			}

		}
	}

	return o
}
