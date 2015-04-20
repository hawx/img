package pixelate

import (
	"image"
	"image/color"

	"hawx.me/code/img/utils"
)

// Halves the width of the double-width image created by hxl to produce nice
// smooth edges.
func halveWidth(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx()/2, b.Dy()))

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx()/2; x++ {
			l := img.At(x*2, y)
			r := img.At(x*2+1, y)

			o.Set(x, y, utils.Average(l, r))
		}
	}

	return o
}

// Hxl pixelates the Image into equilateral triangles with the width
// given. These are arranged into hexagonal shapes.
func Hxl(img image.Image, width int) image.Image {
	b := img.Bounds()

	pixelHeight := width * 2
	pixelWidth := width

	cols := b.Dx() / pixelWidth
	rows := b.Dy() / pixelHeight

	o := image.NewRGBA(image.Rect(0, 0, pixelWidth*cols*2, pixelHeight*rows))

	// Note: "Top" doesn't mean above the x-axis, it means in the triangle
	// pointing towards the x-axis.
	inTop := func(x, y float64) bool {
		return (x >= 0 && y >= x) || (x <= 0 && y >= -x)
	}

	// Same for "Bottom" this is the triangle below and pointing towards the
	// x-axis.
	inBottom := func(x, y float64) bool {
		return (x >= 0 && y <= -x) || (x <= 0 && y <= x)
	}

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			north := []color.Color{}
			south := []color.Color{}

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {
					realY := row*pixelHeight + y
					realX := col*pixelWidth + x
					pixel := img.At(realX, realY)

					yOrigin := float64(y - pixelHeight/2)
					xOrigin := float64(x - pixelWidth/2)

					if inTop(xOrigin, yOrigin) {
						north = append(north, pixel)
					} else if inBottom(xOrigin, yOrigin) {
						south = append(south, pixel)
					}
				}
			}

			top := utils.Average(north...)
			bot := utils.Average(south...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth*2; x++ {
					realY := row*pixelHeight + y
					realX := col*pixelWidth*2 + x

					yOrigin := float64(y - pixelHeight/2)
					xOrigin := float64(x - pixelWidth*2/2)

					if inTop(xOrigin, yOrigin) {
						o.Set(realX, realY, top)
					} else if inBottom(xOrigin, yOrigin) {
						o.Set(realX, realY, bot)
					}
				}
			}
		}
	}

	// Now for the shifted version

	offsetY := pixelHeight / 2
	offsetX := pixelWidth / 2

	for col := -1; col < cols; col++ {
		for row := -1; row < rows; row++ {
			north := []color.Color{}
			south := []color.Color{}

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {
					realY := row*pixelHeight + y + offsetY
					realX := col*pixelWidth + x + offsetX

					if realX >= 0 && realX < b.Dx() {
						pixel := img.At(realX, realY)

						yOrigin := float64(y - pixelHeight/2)
						xOrigin := float64(x - pixelWidth/2)

						if inTop(xOrigin, yOrigin) {
							north = append(north, pixel)
						} else if inBottom(xOrigin, yOrigin) {
							south = append(south, pixel)
						}
					}
				}
			}

			top := utils.Average(north...)
			bot := utils.Average(south...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth*2; x++ {
					realY := row*pixelHeight + y + offsetY
					realX := col*pixelWidth*2 + x + offsetX*2

					yOrigin := float64(y - pixelHeight/2)
					xOrigin := float64(x - pixelWidth*2/2)

					if inTop(xOrigin, yOrigin) {
						o.Set(realX, realY, top)
					} else if inBottom(xOrigin, yOrigin) {
						o.Set(realX, realY, bot)
					}
				}
			}
		}
	}

	return halveWidth(o)
}
