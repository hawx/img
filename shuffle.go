package main

import (
	"./shuffle"
	"./utils"
	"fmt"
	"os"
	"flag"
)

var vertical   = flag.Bool("vertical", false, "Use vertical shuffling only")
var horizontal = flag.Bool("horizontal", false, "Use horizontal shuffling only")

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: shuffle [opts]\n" +
			"\n" +
			"  Shuffle takes a png file from STDIN, shuffles the pixels of the image,\n" +
			"  and prints the result to STDOUT. This allows multiple utilities to be\n" +
			"  easily composed.\n" +
			"\n" +
			"  --horizontal     # Use horizontal shuffling only\n" +
			"  --vertical       # Use vertical shuffling only\n" +
			"  --help           # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := utils.ReadStdin()

	if (*vertical && !*horizontal) {
		i = shuffle.Vertically(i)
	} else if (*horizontal && !*vertical) {
		i = shuffle.Horizontally(i)
	} else {
		i = shuffle.Shuffle(i)
	}

	utils.WriteStdout(i)
}
