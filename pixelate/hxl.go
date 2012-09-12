package pixelate

import (
	"../utils"
	"image"
	"image/color"
)

// Halves the width of the double-width image created by hxl to produce nice
// smooth edges.
func halveWidth(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx() / 2, b.Dy()))

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx() / 2; x++ {
			l := img.At(x * 2, y)
			r := img.At(x * 2 + 1, y)

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
	pixelWidth  := width

	cols  := b.Dx() / pixelWidth
	rows  := b.Dy() / pixelHeight

	o := image.NewRGBA(image.Rect(0, 0, pixelWidth * cols * 2, pixelHeight * rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			north := []color.Color{}
			south := []color.Color{}

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x
					pixel := img.At(realX, realY)

					y_origin := float64(y - pixelHeight / 2)
					x_origin := float64(x - pixelWidth / 2)

					if x_origin > 0 && y_origin > x_origin {
						// north-by-north-east
						north = append(north, pixel)
					} else if x_origin > 0 && y_origin < -x_origin {
						// south-by-south-east
						south = append(south, pixel)
					} else if x_origin < 0 && y_origin < x_origin {
						// south-by-south-west
						south = append(south, pixel)
					} else if x_origin < 0 && y_origin > -x_origin {
						// north-by-north-west
						north = append(north, pixel)
					}
				}
			}

			top := utils.Average(north...)
			bot := utils.Average(south...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth * 2; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth * 2 + x

					y_origin := float64(y - pixelHeight / 2)
					x_origin := float64(x - pixelWidth * 2 / 2)

					var toSet color.Color

					if x_origin >= 0 && y_origin >= x_origin {
						// north-by-north-east
						toSet = top
					} else if x_origin >= 0 && y_origin <= -x_origin {
						// south-by-south-east
						toSet = bot
					} else if x_origin <= 0 && y_origin <= x_origin {
						// south-by-south-west
						toSet = bot
					} else if x_origin <= 0 && y_origin >= -x_origin {
						// north-by-north-west
						toSet = top
					}

					if toSet != nil {
						o.Set(realX, realY, toSet)
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

					realY := row * pixelHeight + y + offsetY
					realX := col * pixelWidth + x + offsetX

					if realX >= 0 && realX < b.Dx() {
						pixel := img.At(realX, realY)

						y_origin := float64(y - pixelHeight / 2)
						x_origin := float64(x - pixelWidth / 2)

						if x_origin > 0 && y_origin > x_origin {
							// north-by-north-east
							north = append(north, pixel)
						} else if x_origin > 0 && y_origin < -x_origin {
							// south-by-south-east
							south = append(south, pixel)
						} else if x_origin < 0 && y_origin < x_origin {
							// south-by-south-west
							south = append(south, pixel)
						} else if x_origin < 0 && y_origin > -x_origin {
							// north-by-north-west
							north = append(north, pixel)
						}
					}
				}
			}

			top := utils.Average(north...)
			bot := utils.Average(south...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth * 2; x++ {

					realY := row * pixelHeight + y + offsetY
					realX := col * pixelWidth * 2 + x + offsetX * 2

					y_origin := float64(y - pixelHeight / 2)
					x_origin := float64(x - pixelWidth * 2 / 2)

					var toSet color.Color

					if x_origin >= 0 && y_origin >= x_origin {
						// north-by-north-east
						toSet = top
					} else if x_origin >= 0 && y_origin <= -x_origin {
						// south-by-south-east
						toSet = bot
					} else if x_origin <= 0 && y_origin <= x_origin {
						// south-by-south-west
						toSet = bot
					} else if x_origin <= 0 && y_origin >= -x_origin {
						// north-by-north-west
						toSet = top
					}

					// This needs to set the left and right unpainted zones!
					if toSet != nil {
						o.Set(realX, realY, toSet)
					}
				}
			}
		}
	}

	return halveWidth(o)
}
