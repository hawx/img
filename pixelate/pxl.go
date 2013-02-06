package pixelate

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
)

const (
	// Triangle types for Pxl
	BOTH = iota  // Decide base on closeness of colours in each quadrant
	LEFT         // Create only left triangles
	RIGHT        // Create only right triangles
)

func halve(img image.Image, size utils.Dimension) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx() / 2, b.Dy() / 2))

	for y := 0; y < b.Dy() / 2; y++ {
		for x := 0; x < b.Dx() / 2; x++ {
			tl := img.At(x * 2, y * 2)
			tr := img.At(x * 2 + 1, y * 2)
			br := img.At(x * 2 + 1, y * 2 - 1)
			bl := img.At(x * 2, y * 2 - 1)

			if y % size.H == 0 {
				o.Set(x, y, tl)
			} else if x % size.W == 0 {
				o.Set(x, y, tl)
			} else {
				o.Set(x, y, utils.Average(tl, tr, bl, br))
			}
		}
	}

	return o
}

func pxlDo(img image.Image, triangle int, size utils.Dimension, scaleFactor int) image.Image {

	b := img.Bounds()

	cols  := b.Dx() / size.W
	rows  := b.Dy() / size.H
	ratio := float64(size.H) / float64(size.W)

	o := image.NewRGBA(image.Rect(0, 0, size.W * cols * scaleFactor, size.H * rows * scaleFactor))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			to := []color.Color{}
			ri := []color.Color{}
			bo := []color.Color{}
			le := []color.Color{}

			tc := 0; rc := 0; bc := 0; lc := 0

			for y := 0; y < size.H; y++ {
				for x := 0; x < size.W; x++ {

					realY := row * size.H + y
					realX := col * size.W + x

					y_origin := float64(y - size.H / 2)
					x_origin := float64(x - size.W / 2)

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

				for y := 0; y < size.H * scaleFactor; y++ {
					for x := 0; x < size.W * scaleFactor; x++ {

						realY := row * size.H * scaleFactor + y
						realX := col * size.W * scaleFactor + x

						y_origin := float64(y - size.H * scaleFactor / 2)
						x_origin := float64(x - size.W * scaleFactor / 2)

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

				for y := 0; y < size.H * scaleFactor; y++ {
					for x := 0; x < size.W * scaleFactor; x++ {

						realY := row * size.H * scaleFactor + y
						realX := col * size.W * scaleFactor + x

						y_origin := float64(y - size.H * scaleFactor / 2)
						x_origin := float64(x - size.W * scaleFactor / 2)

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

	return o
}


// Pxl pixelates an Image into right-angled triangles with the dimensions
// given. The triangle direction can be determined by passing the required value
// as triangle; either BOTH, LEFT or RIGHT.
func Pxl(img image.Image, size utils.Dimension, triangle int) image.Image {
	return halve(pxlDo(img, triangle, size, 2), size)
}

// AliasedPxl does the same as Pxl, but does not smooth diagonal edges of the
// triangles. It is faster, but will produce bad results if size is non-square.
func AliasedPxl(img image.Image, size utils.Dimension, triangle int) image.Image {
	return pxlDo(img, triangle, size, 1)
}
