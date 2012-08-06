package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"fmt"
	"flag"
)

func brightness(img image.Image, value float64) image.Image {
	value  = (100 + value) / 100
	value *= value

	f := func(c color.Color) color.Color {
		r, g, b, a := utils.RatioRGBA(c)

		r = utils.Truncatef(r * value * 255)
		g = utils.Truncatef(g * value * 255)
		b = utils.Truncatef(b * value * 255)

		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a * 255)}
	}

	return utils.ChangePixels(img, f)
}

func printHelp() {
	msg := "Usage: brightness [options]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the brightness using the value given\n" +
		"  and prints the result to STDOUT.\n" +
		"\n" +
		"  --by          # Amount to adjust brightness by\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}


var help = flag.Bool("help", false, "Display this help message")
var by   = flag.Float64("by", 20.0, "Amount to shift brightness by")

func main() {
	flag.Parse()

	if *help {
		printHelp()
	}

	i := utils.ReadStdin()
	i  = brightness(i, *by)
	utils.WriteStdout(i)
}
