// Package utils provides useful helper functions for working with Images and
// Colors.
package utils

import (
	"flag"
	"fmt"
	"os"

	"image"

	"image/png"
	"image/jpeg"
	_ "image/gif"
	"code.google.com/p/go.image/tiff"
)

type output int
const (
	PNG output = iota
	JPEG
	TIFF
)

var Output output = PNG

// ReadStdin reads an image file (either PNG, JPEG or GIF) from standard input.
func ReadStdin() image.Image {
	img, _, _ := image.Decode(os.Stdin)
	return img
}

// WriteStdout writes an Image to standard output as a PNG file.
func WriteStdout(img image.Image) {
	switch Output {
	case JPEG:
		jpeg.Encode(os.Stdout, img, nil)
	case PNG:
		png.Encode(os.Stdout, img)
	case TIFF:
		tiff.Encode(os.Stdout, img, nil)
	}
}

// Warn prints a message to standard error
func Warn(s... interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}

// FlagVisited determines whether the named flag has been visited in the FlagSet
// given. This is helpful if you want to have a flag that triggers an action
// when given, but is not a boolean flag.
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
