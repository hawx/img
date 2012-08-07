package main

import (
	"./utils"
	"image"
	"image/color"
	"fmt"
	"os"
	"flag"
	"math"
)

func saturationShift(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToHSLA(c)
		s += amount
		return utils.ToRGBA(h,s,l,a)
	}

	return utils.ChangePixels(img, f)
}

func simplifiedSaturationShift(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToSimplifiedHSLA(c)
		s += math.Pi * amount
		return utils.ToSimplifiedRGBA(h,s,l,a)
	}

	return utils.ChangePixels(img, f)
}

func printHelp() {
	msg := "Usage: saturation [options]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the saturation by the value given and\n" +
		"  prints the result to STDOUT.\n" +
		"\n" +
		"  --by             # Amount to shift hue by\n" +
		"  --simplified     # Use simplified colour conversion algorithms\n" +
		"  --help           # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help       = flag.Bool("help", false, "Display this help message")
var simplified = flag.Bool("simplified", false, "Use simplified colour conversion algorithms")
var by         = flag.Float64("by", 0.1, "Amount to shift saturation by")

func main() {
	flag.Parse()
	f := saturationShift

	if *help {
		printHelp()
	}
	if *simplified {
		f = simplifiedSaturationShift
	}

	i := utils.ReadStdin()
	i  = f(i, *by)
	utils.WriteStdout(i)
}
