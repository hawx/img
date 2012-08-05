package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"fmt"
	"strconv"
)

func adjustContrast(img image.Image, value float64) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	value = (100 + value) / 100
	value = value * value

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := utils.RatioRGBA(img.At(x, y))

			r = (((r - 0.5) * value) + 0.5) * 255
			g = (((g - 0.5) * value) + 0.5) * 255
			b = (((b - 0.5) * value) + 0.5) * 255
			a = a * 255

			r = utils.TruncateFloat(r)
			g = utils.TruncateFloat(g)
			b = utils.TruncateFloat(b)

			o.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	return o
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

func main() {
	value := float64(15.0)
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			printHelp()
		}

		a, _ := strconv.ParseFloat(os.Args[1], 64)
		value = float64(a)
	}

	i := utils.ReadStdin()
	i  = adjustContrast(i, value)
	utils.WriteStdout(i)
}
