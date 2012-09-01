package main

import (
	"./brightness"
	"./utils"
	"os"
	"fmt"
	"flag"
)

func printHelp() {
	msg := "Usage: brightness [options]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the brightness using the value given\n" +
		"  and prints the result to STDOUT.\n" +
		"\n" +
		"  --by          # Amount to adjust brightness by (default: 20.0)\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}


var help = flag.Bool("help", false, "Display this help message")
var by   = flag.Float64("by", 20.0, "Amount to shift brightness by")

func main() {
	flag.Parse()

	if *help {
		printHelp()
	}

	i := utils.ReadStdin()
	i  = brightness.Adjust(i, *by)
	utils.WriteStdout(i)
}
