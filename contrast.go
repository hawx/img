package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"fmt"
	"strconv"
)

func adjustContrast(img image.Image, value float32) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	value = (100 + value) / 100
	value = value * value

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			rn := float32(uint8(r)) / 255
			gn := float32(uint8(g)) / 255
			bn := float32(uint8(b)) / 255

			rn  = (((rn - 0.5) * value) + 0.5) * 255
			gn  = (((gn - 0.5) * value) + 0.5) * 255
			bn  = (((bn - 0.5) * value) + 0.5) * 255

			if rn > 255 { rn = 255 } else if rn < 0 { rn = 0 }
			if gn > 255 { gn = 255 } else if gn < 0 { gn = 0 }
			if bn > 255 { bn = 255 } else if bn < 0 { bn = 0 }

			rf := uint8(rn)
			gf := uint8(gn)
			bf := uint8(bn)

			o.Set(x, y, color.RGBA{rf, gf, bf, uint8(a)})
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
	value := float32(15.0)
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			printHelp()
		}

		a, _ := strconv.ParseFloat(os.Args[1], 32)
		value = float32(a)
	}

	i := utils.ReadStdin()
	i  = adjustContrast(i, value)
	utils.WriteStdout(i)
}
