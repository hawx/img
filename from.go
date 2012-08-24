package main

import (
	"./utils"
	"os"
	"path/filepath"
	"flag"
	"fmt"

	"image"
	"image/png"
	"image/jpeg"
	"image/gif"
)

const (
	GIF = iota
	JPG
	PNG
)

var extMap = map[string] int {
	".gif":  GIF,
	".jpg":  JPG,
	".jpeg": JPG,
	".png":  PNG,
}

func printHelp() {
	msg := "Usage: from <input>\n" +
		"\n" +
		"  Takes an input image (either PNG, JPEG or GIF) and prints it as a PNG\n" +
	  "  to STDOUT.\n" +
		"\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()
	if *help { printHelp() }

	path := os.Args[1]
	f, _ := os.Open(path)
	var img image.Image

	switch extMap[filepath.Ext(path)] {
	case PNG: img, _ = png.Decode(f)
	case JPG: img, _ = jpeg.Decode(f)
	case GIF: img, _ = gif.Decode(f)
	}

	utils.WriteStdout(img)
}
