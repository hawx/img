package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"fmt"
	"flag"
	"strconv"
	"strings"
)

func pxl(img image.Image, pixelHeight, pixelWidth int) image.Image {
	b := img.Bounds()

	cols := b.Dx() / pixelWidth
	rows := b.Dy() / pixelHeight

	o := image.NewRGBA(image.Rect(0, 0, cols * pixelWidth, rows * pixelHeight))

	// Go through rows
	for row := 0; row < rows; row++ {

		// Now through columns along the row
		for col := 0; col < cols; col++ {

			values := make([]color.Color, pixelHeight * pixelWidth)
			i := 0

			// Within the box need to iterate over every point and add value to array
			for x := 0; x < pixelWidth; x++ {
				for y:= 0; y < pixelHeight; y++ {

					realX := col * pixelWidth + x
					realY := row * pixelHeight + y

					values[i] = img.At(realX, realY)

					i++
				}
			}

			avg := utils.Average(values...)

			// Now to draw the new color
			for x := 0; x < pixelWidth; x++ {
				for y := 0; y < pixelHeight; y++ {

					realX := col * pixelWidth + x
					realY := row * pixelHeight + y

					o.Set(realX, realY, avg)
				}
			}
		}
	}

	return o
}

type pixel struct {
	height, width int
}

func (p *pixel) String() string {
	return fmt.Sprint(*p)
}

func (p *pixel) Set(value string) error {
	parts := strings.Split(value, "x")

	y, _ := strconv.Atoi(parts[0])
	x, _ := strconv.Atoi(parts[1])

	*p = pixel{y, x}

	return nil
}

var pixelFlag pixel = pixel{50, 10}
var help = flag.Bool("help", false, "Display this help message")

func init() {
	flag.Var(&pixelFlag, "size", "Size of pixel to pixelate with")
}

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: colourpixels [opts]\n" +
			"\n" +
			"  Made to pixelate an image, but there are various overflows in calculations\n" +
			"  which produces strange (but mostly pretty) patterns.\n" +
			"\n" +
			"  --size HxW    # Set size of pixel to use, defaults to 20x20\n" +
			"  --help        # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}


	i := utils.ReadStdin()
	i  = pxl(i, pixelFlag.height, pixelFlag.width)
	utils.WriteStdout(i)
}
