package main

import (
	"./utils"
	"os"
	"path/filepath"
	"flag"
	"fmt"

	"image/png"
	"image/jpeg"
)

var extMap = map[string] int {
	".jpg":  0,
	".jpeg": 0,
	".png":  1,
}

func printHelp() {
	msg := "Usage: to <output>\n" +
		"\n" +
		"  Takes an png image from STDIN and writes it to the file specified in\n" +
		"  the format implied by the filename (.jpg, .jpeg = JPEG; .png = PNG).\n" +
		"\n" +
		"  --quality     # Output image quality\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help = flag.Bool("help", false, "Display this help message")
var qual = flag.Int("quality", 80, "Output image quality")

func main() {
	flag.Parse()
	if *help { printHelp() }

	path := flag.Args()[0]
	img  := utils.ReadStdin()
	f, _ := os.Create(path)

	switch extMap[filepath.Ext(path)] {
	case 1:
		png.Encode(f, img)

	case 0:
		opts  := jpeg.Options{*qual}
		jpeg.Encode(f, img, &opts)
	}
}
