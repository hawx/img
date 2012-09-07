package main

import (
	"./blend"
	"./utils"
	"flag"
	"os"
	"fmt"
	"image"
	"image/png"
)



func printHelp() {
	msg := "Usage: compose <other> [opts]\n" +
		"\n" +
		"  Takes a png file from STDIN and blends it with the png file at <other>\n" +
		"  using the method chosen. The result is printed to STDOUT.\n" +
		"\n" +
		"  Note: 'blend colour' refers to the colour taken from <other>, while 'base \n" +
		"  color' refers to the colour taken from STDIN.\n" +
		"\n" +
		"  --opacity      # Opacity of blended image layer (default: 0.5)\n" +
		"\n" +
		"  BASIC\n" +
		"  --normal       # Paints pixels using <other> (default)\n" +
		"  --dissolve     # Paints pixels from <other> randomly, depending on opacity\n" +
		"\n" +
		"  --help         # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}


var help         = flag.Bool("help", false, "")
var opacity      = flag.Float64("opacity", 1.0, "")

// BASIC
var normalM      = flag.Bool("normal", false, "")
var dissolveM    = flag.Bool("dissolve", false, "")

func main() {
	flag.Parse()
	if *help { printHelp() }

	a := utils.ReadStdin()

	path := flag.Args()[0]
	f, _ := os.Open(path)
	b, _ := png.Decode(f)
	var img image.Image

	b = blend.Fade(b, *opacity)

	if *normalM {
		img = blend.Normal(a, b)
	} else if *dissolveM {
		img = blend.Dissolve(a, b)
	}

	utils.WriteStdout(img)
}
