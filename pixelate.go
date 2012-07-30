package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"strconv"
	"image"
	"image/png"
	"image/color"
)


func readStdin() image.Image {
	img, _ := png.Decode(os.Stdin)

	return img
}

func writeStdout(img image.Image) {
	png.Encode(os.Stdout, img)
}

func average(cs []color.Color) color.Color {
	var red, green, blue, alpha uint32
	red = 0; green = 0; blue = 0; alpha = 0

	for i := 0; i < len(cs); i++ {
		r, g, b, a := cs[i].RGBA()

		// Need to do some crazy type conversions first
		rn := uint32(uint8(r))
		gn := uint32(uint8(g))
		bn := uint32(uint8(b))
		an := uint32(uint8(a))

		red += rn; green += gn; blue += bn; alpha += an
	}

	return color.RGBA{
		uint8(red   / uint32(len(cs))),
		uint8(green / uint32(len(cs))),
		uint8(blue  / uint32(len(cs))),
		uint8(alpha / uint32(len(cs))),
	}
}

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

			avg := average(values)

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

var pixelFlag pixel = pixel{20, 20}
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

	i := readStdin()
	i  = pixelate(i, pixelFlag.height, pixelFlag.width)
	writeStdout(i)
}
