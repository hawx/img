package main

import (
	"./utils"
	"image"
	"fmt"
	"os"
	"strconv"
)


func hueShift(img image.Image, amount float64) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			h, s, l, a := utils.ToHSLA(img.At(x, y))

			h += amount

			o.Set(x, y, utils.ToRGBA(h, s, l, a))
		}
	}

	return o
}

func printHelp() {
	msg := "Usage: hue [value]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the hue by the value given and\n" +
		"  prints the result to STDOUT.\n" +
		"\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

func main() {
	value := float64(1.0)
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			printHelp()
		}

		a, _ := strconv.ParseFloat(os.Args[1], 64)
		value = float64(a)
	}

	i := utils.ReadStdin()
	i  = hueShift(i, value)
	utils.WriteStdout(i)
}
