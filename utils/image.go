package utils

import (
	"image"
)

// ExcessMode specifies how excess space is dealt with for tools that may
// produce results with different dimensions to the input image.
type ExcessMode int

const (
	// Ignore any "left over" space. So the resultant Rectangles may be smaller
	// than the original given.
	IGNORE ExcessMode = iota

	// Add the excess to the right- and bottom-most Rectangles. So the resultant
	// Rectangles will be the same size, but some edge Rectangles may be larger.
	ADD

	// Separate the excess into new Rectangles. So the resultant Rectangles will
	// be the same size, but some edge Rectangles may be smaller.
	SEPARATE
)

// chopRectangle is the "god" function for rectangle chopping. It takes a
// Rectangle, with the number of rows and columns to split into, along with the
// height and width of the rows and columns. And produces a list of Rectangles.
func chopRectangle(rect image.Rectangle, rows, cols, rowHeight, colWidth int, mode ExcessMode) []image.Rectangle {
	width := rect.Dx()
	height := rect.Dy()

	excessWidth := width % (cols * colWidth)
	excessHeight := height % (rows * rowHeight)

	rs := make([]image.Rectangle, cols*rows)
	i := 0

	for col := 0; col < cols; col++ {
		localWidth := 0
		// If in last column, add extra on
		if mode == ADD && cols == col+1 {
			localWidth = excessWidth
		}

		for row := 0; row < rows; row++ {
			localHeight := 0
			// If in last row, add extra on
			if mode == ADD && rows == row+1 {
				localHeight = excessHeight
			}

			rs[i] = image.Rectangle{
				image.Point{col * colWidth, row * rowHeight},
				image.Point{(col+1)*colWidth + localWidth, (row+1)*rowHeight + localHeight},
			}

			i++
		}
	}

	if mode == SEPARATE {

		// Do bottom row
		if excessHeight > 0 {
			for col := 0; col < cols; col++ {
				rs = append(rs, image.Rectangle{
					image.Point{col * colWidth, rows * rowHeight},
					image.Point{(col + 1) * colWidth, rows*rowHeight + excessHeight},
				})
			}
		}

		// Do rightmost column
		if excessWidth > 0 {
			for row := 0; row < rows; row++ {
				rs = append(rs, image.Rectangle{
					image.Point{cols * colWidth, row * rowHeight},
					image.Point{cols*colWidth + excessWidth, (row + 1) * rowHeight},
				})
			}
		}

		// Do bottom-right corner
		if excessHeight > 0 && excessWidth > 0 {
			rs = append(rs, image.Rectangle{
				image.Point{cols * colWidth, rows * rowHeight},
				image.Point{cols*colWidth + excessWidth, rows*rowHeight + excessHeight},
			})
		}
	}

	return rs
}

// ChopRectangle splits a Rectangle into the number of rows and columns given.
func ChopRectangle(rect image.Rectangle, rows, cols int, mode ExcessMode) []image.Rectangle {
	width := rect.Dx()
	height := rect.Dy()

	colWidth := width / cols
	rowHeight := height / rows

	return chopRectangle(rect, rows, cols, rowHeight, colWidth, mode)
}

// ChopRectangleToSizes splits a Rectangle into smaller Rectangles with the size
// given.
func ChopRectangleToSizes(rect image.Rectangle, rowHeight, colWidth int, mode ExcessMode) []image.Rectangle {
	width := rect.Dx()
	height := rect.Dy()

	cols := width / colWidth
	rows := height / rowHeight

	return chopRectangle(rect, rows, cols, rowHeight, colWidth, mode)
}
