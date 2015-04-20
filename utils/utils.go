// Package utils provides useful helper functions for working with Images and
// Colors.
package utils

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/tiff"
	"hawx.me/code/img/exif"
)

// This is a string, and not an int of some kind, so that it is easy to find out
// which one is set AND have the type niceness. These are the canonical
// names. Easy.
type output string

// Formats that can be created by img.
const (
	PNG  output = "png"
	JPEG        = "jpeg"
	TIFF        = "tiff"
)

// Output is the currently selected output type, and shouldn't be modified by
// anything other than GetOutput.
var Output = PNG

// GetOutput takes the os.Args array as input, uses it to set the output format,
// then returns the new os.Args array.
//
//   os.Args = utils.GetOutput(os.Args)
//   flag.Parse()
//
func GetOutput(args []string) []string {
	Output = output(args[1])
	return append(args[:1], args[2:]...)
}

// ReadStdin reads an image file (either PNG, JPEG or GIF) from standard input.
func ReadStdin() (image.Image, *exif.Exif) {
	img, _, _ := image.Decode(os.Stdin)
	os.Stdin.Seek(0, 0)
	data := exif.Decode(os.Stdin)
	return img, data
}

// WriteStdout writes an Image to standard output as a PNG file.
func WriteStdout(img image.Image, data *exif.Exif) {
	switch Output {
	case JPEG:
		// Create a temporary file for exiftool to use
		tmp, _ := ioutil.TempFile("", "img-utils-exif-")
		path := tmp.Name()

		// Encode the jpeg to this temp file
		err := jpeg.Encode(tmp, img, nil)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Write the exif data to the temp file
		err = data.Write(path)
		if err != nil {
			// // This will generally return an error, and it still works, so can
			// // probably be ignored.
			// log.Println(err)
		}

		// Reopen the temp file. Yes, really. This is due to the fact that, even if
		// I seek to 0, problems still occur. So this works, I'm leaving it.
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
			return
		}

		io.Copy(os.Stdout, f)

		// Make sure the temp file is deleted after
		defer func() { os.Remove(path) }()

	case PNG:
		err := png.Encode(os.Stdout, img)
		if err != nil {
			log.Fatal(err)
			return
		}

	case TIFF:
		err := tiff.Encode(os.Stdout, img, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

// Warn prints a message to standard error
func Warn(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}

// FlagVisited determines whether the named flag has been visited in the FlagSet
// given. This is helpful if you want to have a flag that triggers an action
// when given, but is not a boolean flag.
func FlagVisited(name string, flags flag.FlagSet) (didFind bool) {
	toFind := flags.Lookup(name)

	flags.Visit(func(f *flag.Flag) {
		if f == toFind {
			didFind = true
		}
	})

	return
}
