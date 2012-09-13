package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
	"fmt"
	"os"
	"flag"
)

func printHelp() {
	msg := "Usage: hue [options]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the hue by the value given and\n" +
		"  prints the result to STDOUT.\n" +
		"\n" +
		"  --by             # Amount to shift hue by (default: 60.0)\n" +
		"  --help           # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help = flag.Bool("help", false, "Display this help message")
var by   = flag.Float64("by", 60.0, "Amount to shift hue by")

func main() {
	flag.Parse()

	if *help {
		printHelp()
	}

	i := utils.ReadStdin()
	i  = hsla.Hue(i, *by)
	utils.WriteStdout(i)
}
