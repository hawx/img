package main

import (
	"./utils"
	"os"
	"fmt"
	"strconv"
	"flag"
	"image"
	"image/color"
)

func pxl(img image.Image, pixelSize int) image.Image {
	pixelHeight := pixelSize
	pixelWidth  := pixelSize

	b := img.Bounds()

	cols := b.Dx() / pixelWidth
	rows := b.Dy() / pixelHeight

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

					y_origin := y - pixelHeight / 2
					x_origin := x - pixelWidth / 2

					if y_origin > x_origin && y_origin > -x_origin {
						tc++
						to = append(to, img.At(realX, realY))

					} else if y_origin < x_origin && y_origin > -x_origin {
						rc++
						ri = append(ri, img.At(realX, realY))

					} else if y_origin < x_origin && y_origin < -x_origin {
						bc++
						bo = append(bo, img.At(realX, realY))

					} else if y_origin > x_origin && y_origin < -x_origin {
						lc++
						le = append(le, img.At(realX, realY))

					}
				}
			}

			ato := utils.Average(to...)
			ari := utils.Average(ri...)
			abo := utils.Average(bo...)
			ale := utils.Average(le...)

			if utils.Closeness(ato, ari) > utils.Closeness(ato, ale) {

				top_right   := utils.Average(ato, ari)
				bottom_left := utils.Average(abo, ale)

				for y := 0; y < pixelHeight; y++ {
					for x := 0; x < pixelWidth; x++ {

						realY := row * pixelHeight + y
						realX := col * pixelWidth + x

						y_origin := y - pixelHeight / 2
						x_origin := x - pixelWidth / 2

						if y_origin > x_origin {
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

						y_origin := y - pixelHeight / 2
						x_origin := x - pixelWidth / 2

						if y_origin >= -x_origin {
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

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: pxl [size]\n" +
			"\n" +
			"  Pixelate takes a png file from STDIN, and pixelates it into triangles.\n" +
			"  The result is printed to STDOUT.\n" +
			"\n" +
			"  --help        # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	pixelSize := 20
	if len(os.Args) > 1 {
		pixelSize, _ = strconv.Atoi(os.Args[1])
	}

	i := utils.ReadStdin()
	i  = pxl(i, pixelSize)
	utils.WriteStdout(i)
}
