package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"fmt"
	"flag"
)

func adjustContrast(img image.Image, value float64) image.Image {
	value  = (100 + value) / 100
	value *= value

	f := func(c color.Color) color.Color {
		r,g,b,a := utils.RatioRGBA(c)

		r = utils.Truncatef((((r - 0.5) * value) + 0.5) * 255)
		g = utils.Truncatef((((g - 0.5) * value) + 0.5) * 255)
		b = utils.Truncatef((((b - 0.5) * value) + 0.5) * 255)
		a = a * 255

		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

	return utils.ChangePixels(img, f)
}

func printHelp() {
	msg := "Usage: contrast [value]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the contrast using the value given\n" +
		"  and prints the result to STDOUT.\n" +
		"\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help = flag.Bool("help", false, "Display this help message")
var by   = flag.Float64("by", 15.0, "Amount to shift hue by")

func main() {
	flag.Parse()

	if *help {
		printHelp()
	}

	i := utils.ReadStdin()
	i  = adjustContrast(i, *by)
	utils.WriteStdout(i)
}
