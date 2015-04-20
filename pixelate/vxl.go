package pixelate

import (
	"image"
	"image/color"
	"math"

	"hawx.me/code/img/channel"
	"hawx.me/code/img/utils"
)

// Vxl pixelates the Image into isometric cubes. It averages the colours and
// naïvely darkens and lightens the colours to mimic highlight and shade.
func Vxl(img image.Image, height int, flip bool, top, left, right float64) image.Image {
	b := img.Bounds()

	pixelHeight := height
	pixelWidth := int(math.Sqrt(3.0) * float64(pixelHeight) / 2.0)

	cols := b.Dx() / pixelWidth
	rows := b.Dy() / pixelHeight

	// intersection of lines
	c := float64(pixelHeight) / 2
	// gradient of lines
	k := math.Sqrt(3.0) / 3.0
	o := image.NewRGBA(image.Rect(0, 0, pixelWidth*cols, pixelHeight*rows))

	// See: http://www.flickr.com/photos/hawx-/8466236036/
	inTopSquare := func(x, y float64) bool {
		if !flip {
			y *= -1
		}
		return y <= -k*x+c && y >= k*x && y >= -k*x && y <= k*x+c
	}

	inBottomRight := func(x, y float64) bool {
		if !flip {
			y *= -1
		}
		return x >= 0 && y <= k*x && y >= k*x-c
	}

	inBottomLeft := func(x, y float64) bool {
		if !flip {
			y *= -1
		}
		return x <= 0 && y <= -k*x && y >= -k*x-c
	}

	inHexagon := func(x, y float64) bool {
		return inTopSquare(x, y) || inBottomRight(x, y) || inBottomLeft(x, y)
	}

	topL := channel.AdjustC(utils.Multiplier(top), channel.Lightness)
	rightL := channel.AdjustC(utils.Multiplier(right), channel.Lightness)
	leftL := channel.AdjustC(utils.Multiplier(left), channel.Lightness)

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			seen := []color.Color{}

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {
					realY := row*(pixelHeight+int(c)) + y
					realX := col*pixelWidth + x
					pixel := img.At(realX, realY)

					yOrigin := float64(y - pixelHeight/2)
					xOrigin := float64(x - pixelWidth/2)

					if inHexagon(xOrigin, yOrigin) {
						seen = append(seen, pixel)
					}
				}
			}

			average := utils.Average(seen...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {
					realY := row*(pixelHeight+int(c)) + y
					realX := col*pixelWidth + x

					yOrigin := float64(y - pixelHeight/2)
					xOrigin := float64(x - pixelWidth/2)

					// This stops white bits showing above the top squares. It does mean
					// the dimensions aren't perfect, but what did you expect with pixels
					// and trig. It is inefficient though, maybe fix that later?
					if (!flip && yOrigin < 0) || (flip && yOrigin > 0) {
						o.Set(realX, realY, topL(average))
					} else {
						if xOrigin > 0 {
							o.Set(realX, realY, rightL(average))
						} else {
							o.Set(realX, realY, leftL(average))
						}
					}

					if inTopSquare(xOrigin, yOrigin) {
						o.Set(realX, realY, topL(average))
					}
					if inBottomRight(xOrigin, yOrigin) {
						o.Set(realX, realY, rightL(average))
					}
					if inBottomLeft(xOrigin, yOrigin) {
						o.Set(realX, realY, leftL(average))
					}
				}
			}
		}
	}

	offsetY := (pixelHeight + int(c)) / 2.0
	offsetX := pixelWidth / 2

	for col := -1; col < cols; col++ {
		for row := -1; row < rows; row++ {
			seen := []color.Color{}

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {
					realY := row*(pixelHeight+int(c)) + y + offsetY
					realX := col*pixelWidth + x + offsetX

					if image.Pt(realX, realY).In(b) {
						pixel := img.At(realX, realY)

						yOrigin := float64(y - pixelHeight/2)
						xOrigin := float64(x - pixelWidth/2)

						if inHexagon(xOrigin, yOrigin) {
							seen = append(seen, pixel)
						}
					}
				}
			}

			if len(seen) <= 0 {
				continue
			}
			average := utils.Average(seen...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {
					realY := row*(pixelHeight+int(c)) + y + offsetY
					realX := col*pixelWidth + x + offsetX

					yOrigin := float64(y - pixelHeight/2)
					xOrigin := float64(x - pixelWidth/2)

					if inTopSquare(xOrigin, yOrigin) {
						o.Set(realX, realY, topL(average))
					}
					if inBottomRight(xOrigin, yOrigin) {
						o.Set(realX, realY, rightL(average))
					}
					if inBottomLeft(xOrigin, yOrigin) {
						o.Set(realX, realY, leftL(average))
					}
				}
			}
		}
	}

	return o
}
