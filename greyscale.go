package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"flag"
	"fmt"
)

type PixelAlterer func(r, g, b uint32) uint8

func alterPixels(img image.Image, f PixelAlterer) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := utils.NormalisedRGBA(img.At(x, y))

			grey := f(r, g, b)

			o.Set(x, y, color.RGBA{grey, grey, grey, uint8(a)})
		}
	}

	return o
}


// Creates a greyscale version of +img+ using the average method, simply averages
// R, G and B values.
func average(r, g, b uint32) uint8 {
	return uint8((r + g + b) / 3)
}

// Creates a grayscale version of +img+ using the lightness method. That is, it
// averages the most prominent and least prominent colours.
func lightness(r, g, b uint32) uint8 {
	maxi := utils.Max(r, g, b)
	mini := utils.Min(r, g, b)

	return uint8((maxi + mini) / 2)
}

// Creates a greyscale version of +img+ using the luminosity method. This uses a
// weighted average to account for humans sensitivity to green above other colours.
// Formula is R * 0.21 + G * 0.71 + B * 0.07.
func luminosity(r, g, b uint32) uint8 {
	return uint8(float32(r) * 0.21 + float32(g) * 0.71 + float32(b) * 0.07)
}

func maximal(r, g, b uint32) uint8 {
	return uint8(utils.Max(r, g, b))
}

func minimal(r, g, b uint32) uint8 {
	return uint8(utils.Min(r, g, b))
}

// Supposed photoshop luminosity method for greyscale.
func photoshop(r, g, b uint32) uint8 {
	return uint8(float32(r) * 0.299 + float32(g) * 0.587 + float32(b) * 0.114)
}

var averageM    = flag.Bool("average",    false, "Use average method")
var lightnessM  = flag.Bool("lightness",  false, "Use lightness method")
var luminosityM = flag.Bool("luminosity", false, "Use standard luminosity method")
var maximalM    = flag.Bool("maximal",    false, "Use maximal decomposition")
var minimalM    = flag.Bool("minimal",    false, "Use minimal decomposition")
var photoshopM  = flag.Bool("photoshop",  false, "Use photoshop luminosity method (default)")

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: greyscale [opts]\n" +
			"\n" +
			"  Takes a png file from STDIN and prints a greyscale version to STDOUT.\n" +
			"  Can choose one of the many greyscaling methods.\n" +
			"\n" +
			"  --average           # Use average method\n" +
			"  --lightness         # Use lightness method\n" +
			"  --luminosity        # Use standard luminosity method\n" +
			"  --maximal           # Use maximal decompositon\n" +
			"  --minimal           # Use minimal decompositon\n" +
			"  --photoshop         # Use photoshop luminosity method (default)\n" +
			"  --help              # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := utils.ReadStdin()

	if *averageM {
		i = alterPixels(i, average)
  } else if *lightnessM {
    i = alterPixels(i, lightness)
	} else if *luminosityM {
		i = alterPixels(i, luminosity)
	} else if *maximalM {
		i = alterPixels(i, maximal)
	} else if *minimalM {
		i = alterPixels(i, minimal)
  } else {
    i = alterPixels(i, photoshop)
  }

	utils.WriteStdout(i)
}
