package main

import (
	"os"
	"fmt"
	"strconv"
	"flag"
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

func closeness(one, two color.Color) uint32 {
	a, b, c, d := one.RGBA()
	w, x, y, z := two.RGBA()

	na := uint32(uint8(a))
	nb := uint32(uint8(b))
	nc := uint32(uint8(c))
	nd := uint32(uint8(d))

	nw := uint32(uint8(w))
	nx := uint32(uint8(x))
	ny := uint32(uint8(y))
	nz := uint32(uint8(z))

	return (na - nw) + (nb - nx) + (nc - ny) + (nd - nz)
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

			ato := average(to)
			ari := average(ri)
			abo := average(bo)
			ale := average(le)

			if closeness(ato, ari) > closeness(ato, ale) {

				top_right := average([]color.Color{ato, ari})
				bottom_left := average([]color.Color{abo, ale})

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

				top_left := average([]color.Color{ato, ale})
				bottom_right := average([]color.Color{abo, ari})

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

	i := readStdin()
	i  = pxl(i, pixelSize)
	writeStdout(i)
}
