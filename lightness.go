package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
	"fmt"
	"os"
	"flag"
)

func printHelp() {
	msg := "Usage: lightness [options]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the lightness by the value given and\n" +
		"  prints the result to STDOUT.\n" +
		"\n" +
		"  --by             # Amount to shift lightness by (default: 0.1)\n" +
		"  --help           # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help = flag.Bool("help", false, "Display this help message")
var by   = flag.Float64("by", 0.1, "Amount to shift lightness by")

func main() {
	flag.Parse()

	if *help {
		printHelp()
	}

	i := utils.ReadStdin()
	i  = hsla.Lightness(i, *by)
	utils.WriteStdout(i)
}
