package pixelate

import (
	"../utils"
	"image"
	"image/color"
)

// Pixelates +img+ into large averaged pixels of size pixelHeight by pixelWidth.
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
