package main

import (
	"./utils"
	"fmt"
	"os"
	"math/rand"
	"image"
	"flag"
)


func randBetween(a, b int) int {
	return rand.Intn(b - a)
}

func shuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			xb := randBetween(b.Min.X, b.Max.X)
			yb := randBetween(b.Min.Y, b.Max.Y)

			a := img.At(x, y)
			b := img.At(xb, yb)

			o.Set(x, y, b)
			o.Set(xb, yb, a)
		}
	}

	return o
}

func verticalShuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			yb := randBetween(b.Min.Y, b.Max.Y)

			a := img.At(x, y)
			b := img.At(x, yb)

			o.Set(x, y, b)
			o.Set(x, yb, a)
		}
	}

	return o
}

func horizontalShuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			xb := randBetween(b.Min.X, b.Max.X)

			a := img.At(x, y)
			b := img.At(xb, y)

			o.Set(x, y, b)
			o.Set(xb, y, a)
		}
	}

	return o
}

var vertical   = flag.Bool("v", false, "Use vertical shuffling only")
var horizontal = flag.Bool("h", false, "Use horizontal shuffling only")

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: shuffle [opts]\n" +
			"\n" +
			"  Shuffle takes a png file from STDIN, shuffles the pixels of the image,\n" +
			"  and prints the result to STDOUT. This allows multiple utilities to be\n" +
			"  easily composed.\n" +
			"\n" +
			"  -h            # Use horizontal shuffling only\n" +
			"  -v            # Use vertical shuffling only\n" +
			"  --help        # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := utils.ReadStdin()

	if (*vertical && !*horizontal) {
		i = verticalShuffle(i)
	} else if (*horizontal && !*vertical) {
		i = horizontalShuffle(i)
	} else {
		i = shuffle(i)
	}

	utils.WriteStdout(i)
}
