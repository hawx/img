// Package crop provides functions for cropping an image to a smaller size.
package crop

import (
	"image"
	"image/color"
	"math"

	"hawx.me/code/img/utils"
)

// cropTo draws a new Image with pixels that have coordinates that when passed
// to the given function returns true.
func cropTo(img image.Image, in func(x, y int) bool) image.Image {
	return cropToValue(img, func(x, y int) float64 {
		if in(x, y) {
			return 1.0
		}
		return 0.0
	})
}

func cropToValue(img image.Image, in func(x, y int) float64) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(b)

	leastX := b.Max.X
	leastY := b.Max.Y
	mostX := 0
	mostY := 0

	for y := b.Min.Y; y <= b.Max.Y; y++ {
		for x := b.Min.X; x <= b.Max.X; x++ {
			if in(x, y) > 0 {
				if x < leastX {
					leastX = x
				}
				if y < leastY {
					leastY = y
				}
				if x > mostX {
					mostX = x
				}
				if y > mostY {
					mostY = y
				}

				r, g, b, a := utils.NormalisedRGBA(img.At(x, y))
				o.Set(x, y, color.NRGBA{
					uint8(r),
					uint8(g),
					uint8(b),
					uint8(float64(a) * in(x, y)),
				})
			}
		}
	}

	return o.SubImage(image.Rect(leastX, leastY, mostX, mostY))
}

// Square crops an Image to a square. It will use the widest possible width if
// the given width is negative.
func Square(img image.Image, size int, direction utils.Direction) image.Image {
	b := img.Bounds()

	if size < 0 {
		size = b.Dx()
		if b.Dy() < b.Dx() {
			size = b.Dy()
		}
	}

	// Assume centre
	minX := (b.Min.X + b.Dx()/2) - size/2
	maxX := (b.Min.X + b.Dx()/2) + size/2
	minY := (b.Min.Y + b.Dy()/2) - size/2
	maxY := (b.Min.Y + b.Dy()/2) + size/2

	switch direction {
	case utils.TopLeft, utils.Top, utils.TopRight:
		minY = b.Min.Y
		maxY = b.Min.Y + size

	case utils.BottomLeft, utils.Bottom, utils.BottomRight:
		minY = b.Max.Y - size
		maxY = b.Max.Y
	}

	switch direction {
	case utils.TopLeft, utils.Left, utils.BottomLeft:
		minX = b.Min.X
		maxX = b.Min.X + size

	case utils.TopRight, utils.Right, utils.BottomRight:
		minX = b.Max.X - size
		maxX = b.Max.X
	}

	in := func(x, y int) bool {
		return x >= minX && x <= maxX && y >= minY && y <= maxY
	}

	return cropTo(img, in)
}

// Circle crops an Image to a circle. It will use the widest possible width if
// the given width is negative.
func Circle(img image.Image, size int, direction utils.Direction) image.Image {
	b := img.Bounds()
	h := b.Dy()
	w := b.Dx()

	if size < 0 {
		size = w
		if h < w {
			size = h
		}
	}
	size /= 2

	// Assume centre
	xOffset := w / 2
	yOffset := h / 2

	switch direction {
	case utils.TopLeft, utils.Top, utils.TopRight:
		yOffset = size

	case utils.BottomLeft, utils.Bottom, utils.BottomRight:
		yOffset = (h - size)
	}

	switch direction {
	case utils.TopLeft, utils.Left, utils.BottomLeft:
		xOffset = size

	case utils.TopRight, utils.Right, utils.BottomRight:
		xOffset = (w - size)
	}

	in := func(x, y int) float64 {
		x -= xOffset
		y -= yOffset
		if (x*x)+(y*y) < (size)*(size-1) {
			return 1
		} else if (x*x)+(y*y) <= size*size {
			return 0.4
		}
		return 0
	}

	return cropToValue(img, in)
}

// Triangle crops an Image to an equilaterla triangle. It will use the widest
// possible width if the given size is negative.
func Triangle(img image.Image, size int, direction utils.Direction) image.Image {
	b := img.Bounds()
	k := math.Sqrt(3)

	if size < 0 {
		size = b.Dx()
	}
	height := int(k / 2 * float64(size))

	// Assume centre
	minX := (b.Min.X + b.Dx()/2) - size/2
	minY := (b.Min.Y + b.Dy()/2) - height/2
	maxY := (b.Min.Y + b.Dy()/2) + height/2

	switch direction {
	case utils.TopLeft, utils.Top, utils.TopRight:
		minY = b.Min.Y
		maxY = b.Min.Y + height

	case utils.BottomLeft, utils.Bottom, utils.BottomRight:
		minY = b.Max.Y - height
		maxY = b.Max.Y
	}

	switch direction {
	case utils.TopLeft, utils.Left, utils.BottomLeft:
		minX = b.Min.X

	case utils.TopRight, utils.Right, utils.BottomRight:
		minX = b.Max.X - size
	}

	in := func(x, y int) bool {
		x -= minX

		return y <= maxY && y >= minY &&
			y >= int(k*float64(x))-2*height+maxY &&
			y >= int(k*float64(-x))+maxY
	}

	return cropTo(img, in)
}
