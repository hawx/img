package main

import (
	"./utils"
	"os"
	"fmt"
	"flag"
	"image"
	"image/color"
)

func pixelate(img image.Image, pixelHeight, pixelWidth int) image.Image {
	b := img.Bounds()

	cols := b.Dx() / pixelWidth
	rows := b.Dy() / pixelHeight

	o := image.NewRGBA(image.Rect(0, 0, pixelWidth * cols, pixelHeight * rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			values := make([]color.Color, pixelHeight * pixelWidth)
			count  := 0

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x

					values[count] = img.At(realX, realY)
					count++
				}
			}

			avg := utils.Average(values...)

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x

					o.Set(realX, realY, avg)
				}
			}

		}
	}

	return o
}


var pixelFlag utils.Pixel = utils.Pixel{20, 20}
var help = flag.Bool("help", false, "Display this help message")

func init() {
	flag.Var(&pixelFlag, "size", "Size of pixel to pixelate with")
}

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: pixelate [opts]\n" +
			"\n" +
			"  Pixelate takes a png file from STDIN, and pixelates it by averaging the,\n" +
			"  colors in large areas. The result is printed to STDOUT.\n" +
			"\n" +
			"  --size HxW    # Set size of pixel to use, defaults to 20x20\n" +
			"  --help        # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := utils.ReadStdin()
	i  = pixelate(i, pixelFlag.H, pixelFlag.W)
	utils.WriteStdout(i)
}
