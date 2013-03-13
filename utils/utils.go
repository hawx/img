// Package utils provides useful helper functions for working with Images and
// Colors.
package utils

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"image"
	"image/color"
	"image/draw"

	"image/png"
	_ "image/jpeg"
	_ "image/gif"
)

// ReadStdin reads an image file from standard input.
func ReadStdin() image.Image {
	img, _, _ := image.Decode(os.Stdin)

	return img
}

// WriteStdout writes an Image to standard output.
func WriteStdout(img image.Image) {
	png.Encode(os.Stdout, img)
}

// Warn prints a message to standard error
func Warn(s... interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}

// NormalisedRGBA returns the RGBA colour channel values of a Color in a
// normalised form.
func NormalisedRGBA(c color.Color) (rn, gn, bn, an uint32) {
	d := color.NRGBAModel.Convert(c).(color.NRGBA)
	r := d.R; g := d.G; b := d.B; a := d.A

	// Need to do some crazy type conversions first
	rn = uint32(uint8(r))
	gn = uint32(uint8(g))
	bn = uint32(uint8(b))
	an = uint32(uint8(a))

	return
}

// NormalisedRGBAf returns the RGBA colour channel values as floating point
// numbers with values from 0 to 255.
func NormalisedRGBAf(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := NormalisedRGBA(c)

	rn = float64(r)
	gn = float64(g)
	bn = float64(b)
	an = float64(a)

	return
}

// RatioRGBA returns the RGBA colour channel values as floating point numbers
// with values from 0 to 1.
func RatioRGBA(c color.Color) (rn, gn, bn, an float64) {
	r, g, b, a := NormalisedRGBAf(c)

	rn = r / 255
	gn = g / 255
	bn = b / 255
	an = a / 255

	return
}

// Truncate takes a colour channel value and forces it into the range 0 to 255
// by setting any value below 0 to 0 and and any above 255 to 255.
func Truncate(n uint32) uint32 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

// Truanctef is identical to Truncate but takes and returns a float64.
func Truncatef(n float64) float64 {
	if n < 0 { return 0 } else if n > 255 { return 255 }
	return n
}

// Closeness calculates the "closeness" of two colours by finding the sum of
// differences in each colour channel.
func Closeness(one, two color.Color) uint32 {
	a, b, c, d := NormalisedRGBA(one)
	w, x, y, z := NormalisedRGBA(two)

	return (a - w) + (b - x) + (c - y) + (d - z)
}

// Average takes a list of colours and returns the average.
func Average(cs... color.Color) color.Color {
	if len(cs) < 1 {
		return color.Black
	}

	var red, green, blue, alpha uint32
	red = 0; green = 0; blue = 0; alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := NormalisedRGBA(cs[i])

		red += r; green += g; blue += b; alpha += a
	}

	return color.NRGBA{
		uint8(red   / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue  / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}

// Min returns the smallest value in the list of uint32s given.
func Min(ns... uint32) (n uint32) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] < n {
			n = ns[i]
		}
	}
	return
}

// Minf returns the smallest value in the list of float64s given.
func Minf(ns... float64) (n float64) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] < n {
			n = ns[i]
		}
	}
	return
}

// Max returns the largest value in the list of uint32s given.
func Max(ns... uint32) (n uint32) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] > n {
			n = ns[i]
		}
	}
	return
}

// Maxf returns the largest value in the list of float64s given.
func Maxf(ns... float64) (n float64) {
	if len(ns) > 0 {
		n = ns[0]
	}
	for i := 1; i < len(ns); i++ {
		if ns[i] > n {
			n = ns[i]
		}
	}
	return
}


func ChopRectangle(rect image.Rectangle, cols, rows int) []image.Rectangle {
	w := rect.Dx()
	h := rect.Dy()

	colWidth  := w / cols
	rowHeight := h / rows

	rs := make([]image.Rectangle, cols * rows)
	i  := 0

	for col := 0; col < cols; col++ {
		// If in last column, add extra on
		if cols == col - 1 {
			colWidth += w % cols
		}

		for row := 0; row < rows; row++ {
			// If in last row, add extra on
			if rows == row - 1 {
				rowHeight += h % rows
			}

			rs[i] = image.Rectangle{
				image.Point{col     * colWidth, row     * rowHeight},
				image.Point{(col+1) * colWidth, (row+1) * rowHeight},
			}

			i++
		}
	}

	return rs
}

func ChopRectangleToSizesExcess(rect image.Rectangle, height, width int) []image.Rectangle {
	w := rect.Dx()
	h := rect.Dy()

	cols := w / width
	rows := h / height

	excessHeight := h % (rows * height)
	excessWidth  := w % (cols * height)

	rs := ChopRectangleToSizes(rect, height, width)

	// Do bottom row
	if excessHeight > 0 {
		for col := 0; col < cols; col++ {
			rs = append(rs, image.Rectangle{
				image.Point{col     * width, rows * height},
				image.Point{(col+1) * width, rows * height + excessHeight},
			})
		}
	}

	// Do rightmost column
	if excessWidth > 0 {
		for row := 0; row < rows; row++ {
			rs = append(rs, image.Rectangle{
				image.Point{cols * width,               row * height},
				image.Point{cols * width + excessWidth, (row+1) * height},
			})
		}
	}

	// Do bottom-right corner
	if excessHeight > 0 && excessWidth > 0 {
		rs = append(rs, image.Rectangle{
			image.Point{cols * width,               rows * height},
			image.Point{cols * width + excessWidth, rows * height + excessHeight},
		})
	}

	return rs
}

// Ignores excesss.
func ChopRectangleToSizes(rect image.Rectangle, height, width int) []image.Rectangle {
	w := rect.Dx()
	h := rect.Dy()

	cols := w / width
	rows := h / height

	rs := make([]image.Rectangle, cols * rows)
	i  := 0

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			rs[i] = image.Rectangle{
				image.Point{col     * width, row     * height},
				image.Point{(col+1) * width, (row+1) * height},
			}

			i++
		}
	}

	return rs
}

func splitRectangle(b image.Rectangle, parts int) []image.Rectangle {
	if b.Dx() > b.Dy() {
		return ChopRectangle(b, parts, 1)
	}
	return ChopRectangle(b, 1, parts)
}


// EachColor iterates through each pixel of the Image, applying the function
// to each colour.
func EachColor(img image.Image, f func(c color.Color)) {
	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			f(img.At(x, y))
		}
	}
}
// PEachColor is like EachColor, but runs in parallel. This means that order can
// not be guaranteed.
func PEachColor(img image.Image, f func(c color.Color)) {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	c := make(chan int, nCPU)

	for _, r := range splitRectangle(img.Bounds(), nCPU) {
		go EachColorInRectangle(img, r, f)
	}

	// wait until work is done
	for i := 0; i < nCPU; i++ {
		<-c
	}
}

func peachColorWorker(img image.Image, b image.Rectangle, f func(c color.Color), c chan int) {

	EachColorInRectangle(img, b, f)
	c <- 1
}

// EachColorInRectangle is a helper function for working on a part of an image.
func EachColorInRectangle(img image.Image, b image.Rectangle, f func(c color.Color)) {
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			f(img.At(x, y))
		}
	}
}


// MapColor iterates through each pixel of the Image and applies the given
// function, drawing the returned colour to a new Image which is then returned.
func MapColor(img image.Image, f Composable) image.Image {
	// Use maximum number of CPUs available
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	c := make(chan int, nCPU)
	o := image.NewRGBA(img.Bounds())

	for _, r := range splitRectangle(img.Bounds(), nCPU) {
		go mapColorWorker(img, r, o, f, c)
	}

	// wait until work is done
	for i := 0; i < nCPU; i++ {
		<-c
	}

	return o
}

func mapColorWorker(img image.Image, bounds image.Rectangle, dest draw.Image,
	f Composable, c chan int) {

	MapColorInRectangle(img, bounds, dest, f)
	c <- 1
}

// MapColorInRectangle is a helper function for working on part of an image. It
// takes the original image, a function to use, a image to write to, and the
// bounds of the original (and therefore the final image) to act upon.
func MapColorInRectangle(img image.Image, bounds image.Rectangle, dest draw.Image,
	f Composable) {

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dest.Set(x, y, f(img.At(x, y)))
		}
	}
}


func FlagVisited(name string, flags flag.FlagSet) bool {
	didFind := false
	toFind  := flags.Lookup(name)

	flags.Visit(func (f *flag.Flag) {
		if f == toFind {
			didFind = true
		}
	})

	return didFind
}

// Adjuster is a function which takes a floating point value, most likely
// between 0 and 1, and returns another floating point value, again between 0
// and 1 (though this is not a requirement).
type Adjuster func (float64) float64

// Adder returns an Adjuster which adds the value with to the returned Adjusters
// argument.
func Adder(with float64) Adjuster {
	return func(i float64) float64 {
		return i + with
	}
}

// Multiplier returns an Adjuster which multiplies a given value with ratio.
func Multiplier(ratio float64) Adjuster {
	return func(i float64) float64 {
		return i * ratio
	}
}

// Composable is a function which can be composed with another function of the
// same type. It transforms a Color into another Color.
type Composable func(color.Color) color.Color

// Map takes a Composable function, and returns a function which acts on Images.
func Map(f func() Composable) func(image.Image) image.Image {
	return func(img image.Image) image.Image {
		return MapColor(img, f())
	}
}

// MapAdjuster takes a Composable function and an Adjuster, and returns a
// function which acts on Images.
func MapAdjuster(f func(Adjuster) Composable) func(image.Image, Adjuster) image.Image {
	return func(img image.Image, adj Adjuster) image.Image {
		return MapColor(img, f(adj))
	}
}

// Compose takes a variable list of Composable functions and returns a single
// Composable function which performs each of them sequentially. For example,
//
//   var f Composable = Compose(
//     greyscale.BlueC(),
//     brightness.AdjustC(utils.Adder(0.05)),
//   )
//
//   // Only loops through image once!
//   img = MapColors(img, f)
//
func Compose(fs... Composable) Composable {
	return func(c color.Color) color.Color {
		for _, f := range fs {
			c = f(c)
		}
		return c
	}
}
