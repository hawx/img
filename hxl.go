package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
	"os"
	"fmt"
	"flag"
)

var width = flag.Int("width", 20, "Width of the base of each triangle")
var help  = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: hxl [opts]\n" +
			"\n" +
			"  Hxl takes a png file from STDIN, and pixelates it into triangles.\n" +
			"  The result is printed to STDOUT.\n" +
			"\n" +
			"  --width <px>    # Width of the base of each triangle (default: 20)\n" +
			"  --help          # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := utils.ReadStdin()
	i  = pixelate.Hxl(i, *width)
	utils.WriteStdout(i)
}
