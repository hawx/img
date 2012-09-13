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
func Pixelate(img image.Image, pixelHeight, pixelWidth int) image.Image {
	b := img.Bounds()

	cols := b.Dx() / pixelWidth
	rows := b.Dy() / pixelHeight

	o := image.NewRGBA(image.Rect(0, 0, pixelWidth * cols, pixelHeight * rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			values := make([]color.Color, pixelHeight * pixelWidth)
			count  := 0

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x

					values[count] = img.At(realX, realY)
					count++
				}
			}

			avg := utils.Average(values...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x

					o.Set(realX, realY, avg)
				}
			}

		}
	}

	return o
}
