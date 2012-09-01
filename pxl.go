package main

import (
	"./pixelate"
	"./utils"
	"os"
	"fmt"
	"flag"
)

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

	triangle := pixelate.BOTH
	if *left  { triangle = pixelate.LEFT }
	if *right { triangle = pixelate.RIGHT}

	i := utils.ReadStdin()
	i  = pixelate.Pxl(i, triangle, pixelFlag.H, pixelFlag.W)
	utils.WriteStdout(i)
}
