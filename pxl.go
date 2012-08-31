package main

import (
	"./utils"
	"os"
	"fmt"
	"flag"
	"image"
	"image/color"
)

const (
	BOTH = iota
	LEFT
	RIGHT
)

func pxl(img image.Image, triangle, pixelHeight, pixelWidth int) image.Image {
	b := img.Bounds()

	cols  := b.Dx() / pixelWidth
	rows  := b.Dy() / pixelHeight
	ratio := float64(pixelHeight) / float64(pixelWidth)

	o := image.NewRGBA(image.Rect(0, 0, pixelWidth * cols, pixelHeight * rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {

			to := []color.Color{}
			ri := []color.Color{}
			bo := []color.Color{}
			le := []color.Color{}

			tc := 0; rc := 0; bc := 0; lc := 0

			for y := 0; y < pixelHeight; y++ {
				for x := 0; x < pixelWidth; x++ {

					realY := row * pixelHeight + y
					realX := col * pixelWidth + x

					y_origin := float64(y - pixelHeight / 2)
					x_origin := float64(x - pixelWidth / 2)

					if y_origin > ratio * x_origin && y_origin > ratio * -x_origin {
						tc++
						to = append(to, img.At(realX, realY))

					} else if y_origin < ratio * x_origin && y_origin > ratio * -x_origin {
						rc++
						ri = append(ri, img.At(realX, realY))

					} else if y_origin < ratio * x_origin && y_origin < ratio * -x_origin {
						bc++
						bo = append(bo, img.At(realX, realY))

					} else if y_origin > ratio * x_origin && y_origin < ratio * -x_origin {
						lc++
						le = append(le, img.At(realX, realY))

					}
				}
			}

			ato := utils.Average(to...)
			ari := utils.Average(ri...)
			abo := utils.Average(bo...)
			ale := utils.Average(le...)

			if (triangle != LEFT) && (triangle == RIGHT ||
				utils.Closeness(ato, ari) > utils.Closeness(ato, ale)) {

				top_right   := utils.Average(ato, ari)
				bottom_left := utils.Average(abo, ale)

				for y := 0; y < pixelHeight; y++ {
					for x := 0; x < pixelWidth; x++ {

						realY := row * pixelHeight + y
						realX := col * pixelWidth + x

						y_origin := float64(y - pixelHeight / 2)
						x_origin := float64(x - pixelWidth / 2)

						if y_origin > ratio * x_origin {
							o.Set(realX, realY, top_right)
						} else {
							o.Set(realX, realY, bottom_left)
						}
					}
				}

			} else {

				top_left     := utils.Average(ato, ale)
				bottom_right := utils.Average(abo, ari)

				for y := 0; y < pixelHeight; y++ {
					for x := 0; x < pixelWidth; x++ {

						realY := row * pixelHeight + y
						realX := col * pixelWidth + x

						y_origin := float64(y - pixelHeight / 2)
						x_origin := float64(x - pixelWidth / 2)

						if y_origin >= ratio * -x_origin {
							o.Set(realX, realY, top_left)
						} else {
							o.Set(realX, realY, bottom_right)
						}
					}
				}

			}


		}
	}

	return o
}


var pixelFlag utils.Pixel = utils.Pixel{20, 20}
var left  = flag.Bool("left", false, "Create only top-left/bottom-right triangles")
var right = flag.Bool("right", false, "Create only top-right/bottom-left triangles")
var both  = flag.Bool("both", false, "Create traingles based on closeness of colours (default)")
var help  = flag.Bool("help", false, "Display this help message")

func init() {
	flag.Var(&pixelFlag, "size", "Size of pixel to use, defaults to 20x20")
}

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: pxl [opts]\n" +
			"\n" +
			"  Pxl takes a png file from STDIN, and pixelates it into triangles.\n" +
			"  The result is printed to STDOUT.\n" +
			"\n" +
			"  --both          # Create triangles based on closeness of colours (default)\n" +
			"  --left          # Create only top-left/bottom-right triangles\n" +
			"  --right         # Create only top-right/bottom-left triangles\n" +
			"  --size HxW      # Size of pixel to use, defaults to 20x20\n" +
			"  --help          # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	triangle := BOTH
	if *left  { triangle = LEFT }
	if *right { triangle = RIGHT}

	i := utils.ReadStdin()
	i  = pxl(i, triangle, pixelFlag.H, pixelFlag.W)
	utils.WriteStdout(i)
}
