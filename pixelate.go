package main

import (
	"./pixelate"
	"./utils"
	"os"
	"fmt"
	"flag"
)

var pixelFlag utils.Pixel = utils.Pixel{20, 20}
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

	i := utils.ReadStdin()
	i  = pixelate.Pixelate(i, pixelFlag.H, pixelFlag.W)
	utils.WriteStdout(i)
}
