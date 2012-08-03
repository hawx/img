package main

import (
	"os"
	"image"
	"image/png"
	"image/color"
	"flag"
	"fmt"
)


func readStdin() image.Image {
	img, _ := png.Decode(os.Stdin)

	return img
}

func writeStdout(img image.Image) {
	png.Encode(os.Stdout, img)
}

func min(ns... uint32) (n uint32) {
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

func max(ns... uint32) (n uint32) {
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

type PixelAlterer func(r, g, b uint32) uint8

func alterPixels(img image.Image, f PixelAlterer) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			rn := uint32(uint8(r))
			gn := uint32(uint8(g))
			bn := uint32(uint8(b))

			grey := f(rn, gn, bn)

			o.Set(x, y, color.RGBA{grey, grey, grey, uint8(a)})
		}
	}

	return o
}

/*
 Creates a greyscale version of +img+ using the average method, simply averages
 R, G and B values.
 */
	func average(r, g, b uint32) uint8 {
	return uint8((r + g + b) / 3)
}

/*
 Creates a grayscale version of +img+ using the lightness method. That is, it
 averages the most prominent and least prominent colours.
 */
	func lightness(r, g, b uint32) uint8 {
	maxi := max(r, g, b)
	mini := min(r, g, b)

	return uint8((maxi + mini) / 2)
}

/*
 Creates a greyscale version of +img+ using the luminosity method. This uses a
 weighted average to account for humans sensitivity to green above other colours.
 Formula is R * 0.21 + G * 0.71 + B * 0.07.
 */
	func luminosity(r, g, b uint32) uint8 {
	return uint8(float32(r) * 0.21 + float32(g) * 0.71 + float32(b) * 0.07)
}

var averageM    = flag.Bool("average", false, "Use average method")
var lightnessM  = flag.Bool("lightness", false, "Use lightness method")
var luminosityM = flag.Bool("luminosity", false, "Use luminosity method (default)")

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: greyscale [opts]\n" +
			"\n" +
			"  Takes a png file from STDIN and prints a greyscale version to STDOUT.\n" +
			"  Can choose one of three different greyscaling methods.\n" +
			"\n" +
			"  --average           # Use average method\n" +
			"  --lightness         # Use lightness method\n" +
			"  --luminosity        # Use luminosity method (default)\n" +
			"  --help              # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := readStdin()

	if *averageM {
		i = alterPixels(i, average)
  } else if *lightnessM {
    i = alterPixels(i, lightness)
  } else {
    i = alterPixels(i, luminosity)
  }

	writeStdout(i)
}
