package main

import (
	"./contrast"
	"./utils"
	"os"
	"fmt"
	"flag"
)

func printHelp() {
	msg := "Usage: contrast [value]\n" +
		"\n" +
		"  Takes a png file from STDIN, adjusts the contrast using the value given\n" +
		"  and prints the result to STDOUT.\n" +
		"\n" +
		"  --by          # Amount to shift contrast by (default: 15.0)\n" +
		"  --help        # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

var help = flag.Bool("help", false, "Display this help message")
var by   = flag.Float64("by", 15.0, "Amount to shift contrast by")

func main() {
	flag.Parse()

	if *help {
		printHelp()
	}

	i := utils.ReadStdin()
	i  = contrast.Adjust(i, *by)
	utils.WriteStdout(i)
}
