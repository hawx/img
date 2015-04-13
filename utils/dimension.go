package utils

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"
)

// Dimension represents a rectangle with Height and Width.
type Dimension struct {
	H, W int
}

func (d *Dimension) String() string {
	return fmt.Sprint("%vx%v", *d)
}

// Set takes a string representing a Dimension (ie., in the format HxW, where H
// and W are both integers) and sets the Dimension called on to the correct
// value. It returns an error string if a problem occurs.
func (d *Dimension) Set(value string) error {
	parts := strings.Split(value, "x")

	h, e := strconv.Atoi(parts[0])
	if e != nil {
		return errors.New("Error parsing height: expect HxW where H and W are integers")
	}

	w, f := strconv.Atoi(parts[1])
	if f != nil {
		return errors.New("Error parsing width: expect HxW where H and W are integers")
	}

	*d = Dimension{h, w}

	return nil
}

// SizeForRows determines the maximum (square) Dimensions to use that will fit
// the given number of rows into the image.
func SizeForRows(img image.Image, rows int) Dimension {
	b := img.Bounds()
	h := b.Dy() / rows

	return Dimension{h, h}
}

// SizeForCols determines the maximum (square) Dimensions to use that will fit
// the given number of columns into the image.
func SizeForCols(img image.Image, cols int) Dimension {
	b := img.Bounds()
	w := b.Dx() / cols

	return Dimension{w, w}
}

// SizeForRowsAndCols determines the maximum (rectangular) Dimensions to use
// that will fit the given number of rows and columns into the image.
func SizeForRowsAndCols(img image.Image, rows, cols int) Dimension {
	b := img.Bounds()

	return Dimension{b.Dy() / rows, b.Dx() / cols}
}
