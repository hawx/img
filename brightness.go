package main

import (
	"./utils"
	"os"
	"image"
	"image/color"
	"fmt"
	"strconv"
)

func brightness(img image.Image, value float64) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	value = (100 + value) / 100
	value = value * value

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			rn := float64(uint8(r)) / 255
			gn := float64(uint8(g)) / 255
			bn := float64(uint8(b)) / 255

			rn  = (rn * value) * 255
			gn  = (gn * value) * 255
			bn  = (bn * value) * 255

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
	msg := "Usage: brightness [value]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the brightness using the value given\n" +
		"  and prints the result to STDOUT.\n" +
		"\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}


func main() {
	value := float64(20.0)
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			printHelp()
		}

		a, _ := strconv.ParseFloat(os.Args[1], 64)
		value = float64(a)
	}

	i := utils.ReadStdin()
	i  = brightness(i, value)
	utils.WriteStdout(i)
}
