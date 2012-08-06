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

func hueShift(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToHSLA(c)
		h += amount
		return utils.ToRGBA(h,s,l,a)
	}

	return utils.ChangePixels(img, f)
}

func simplifiedHueShift(img image.Image, amount float64) image.Image {
	f := func(c color.Color) color.Color {
		h,s,l,a := utils.ToSimplifiedHSLA(c)
		h += math.Pi * amount
		return utils.ToSimplifiedRGBA(h,s,l,a)
	}

	return utils.ChangePixels(img, f)
}

func printHelp() {
	msg := "Usage: hue [options]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the hue by the value given and\n" +
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
var by         = flag.Float64("by", 60.0, "Amount to shift hue by")

func main() {
	flag.Parse()
	f := hueShift

	if *help {
		printHelp()
	}
	if *simplified {
		f = simplifiedHueShift
	}

	i := utils.ReadStdin()
	i  = f(i, *by)
	utils.WriteStdout(i)
}
