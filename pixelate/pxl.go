package pixelate

import (
	"../utils"
	"image"
	"image/color"
)

const (
	BOTH = iota
	LEFT
	RIGHT
)

func halve(img image.Image, pixelHeight, pixelWidth int) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx() / 2, b.Dy() / 2))

	for y := 0; y < b.Dy() / 2; y++ {
		for x := 0; x < b.Dx() / 2; x++ {
			tl := img.At(x * 2, y * 2)
			tr := img.At(x * 2 + 1, y * 2)
			br := img.At(x * 2 + 1, y * 2 - 1)
			bl := img.At(x * 2, y * 2 - 1)

			if y % pixelHeight == 0 {
				o.Set(x, y, tl)
			} else if x % pixelWidth == 0 {
				o.Set(x, y, tl)
			} else {
				o.Set(x, y, utils.Average(tl, tr, bl, br))
			}
		}
	}

	return o
}

// Pixelates +img+ into right-angled triangles of size +pixelHeight+ by
// +pixelWidth+. +triangle+ determines whether the direction of triangles is
// determined by the closeness of colours (pixelate.BOTH), or only left
// (pixelate.LEFT) or right (pixelate.RIGHT) triangles are made.
func Pxl(img image.Image, triangle, pixelHeight, pixelWidth int) image.Image {
	b := img.Bounds()

	cols  := b.Dx() / pixelWidth
	rows  := b.Dy() / pixelHeight
	ratio := float64(pixelHeight) / float64(pixelWidth)

	o := image.NewRGBA(image.Rect(0, 0, pixelWidth * cols * 2, pixelHeight * rows * 2))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			to := []color.Color{}
			ri := []color.Color{}
			bo := []color.Color{}
			le := []color.Color{}

			tc := 0; rc := 0; bc := 0; lc := 0

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x

					y_origin := float64(y - pixelHeight / 2)
					x_origin := float64(x - pixelWidth / 2)

					if y_origin > ratio * x_origin && y_origin > ratio * -x_origin {
						tc++
						to = append(to, img.At(realX, realY))

					} else if y_origin < ratio * x_origin && y_origin > ratio * -x_origin {
						rc++
						ri = append(ri, img.At(realX, realY))

					} else if y_origin < ratio * x_origin && y_origin < ratio * -x_origin {
						bc++
						bo = append(bo, img.At(realX, realY))

					} else if y_origin > ratio * x_origin && y_origin < ratio * -x_origin {
						lc++
						le = append(le, img.At(realX, realY))

					}
				}
			}

			ato := utils.Average(to...)
			ari := utils.Average(ri...)
			abo := utils.Average(bo...)
			ale := utils.Average(le...)

			if (triangle != LEFT) && (triangle == RIGHT ||
				utils.Closeness(ato, ari) > utils.Closeness(ato, ale)) {

				top_right   := utils.Average(ato, ari)
				bottom_left := utils.Average(abo, ale)

				for y := 0; y < pixelHeight * 2; y++ {
					for x := 0; x < pixelWidth * 2; x++ {

						realY := row * pixelHeight * 2 + y
						realX := col * pixelWidth * 2 + x

						y_origin := float64(y - pixelHeight * 2 / 2)
						x_origin := float64(x - pixelWidth * 2 / 2)

						if y_origin > ratio * x_origin {
							o.Set(realX, realY, top_right)
						} else {
							o.Set(realX, realY, bottom_left)
						}
					}
				}

			} else {

				top_left     := utils.Average(ato, ale)
				bottom_right := utils.Average(abo, ari)

				for y := 0; y < pixelHeight * 2; y++ {
					for x := 0; x < pixelWidth * 2; x++ {

						realY := row * pixelHeight * 2 + y
						realX := col * pixelWidth * 2 + x

						y_origin := float64(y - pixelHeight * 2 / 2)
						x_origin := float64(x - pixelWidth * 2 / 2)

						if y_origin >= ratio * -x_origin {
							o.Set(realX, realY, top_left)
						} else {
							o.Set(realX, realY, bottom_right)
						}
					}
				}

			}


		}
	}

	return halve(o, pixelHeight, pixelWidth)
}
