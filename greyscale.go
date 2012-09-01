package main

import (
	"./greyscale"
	"./utils"
	"os"
	"flag"
	"fmt"
)

var averageM    = flag.Bool("average",    false, "Use average method")
var lightnessM  = flag.Bool("lightness",  false, "Use lightness method")
var luminosityM = flag.Bool("luminosity", false, "Use standard luminosity method")
var maximalM    = flag.Bool("maximal",    false, "Use maximal decomposition")
var minimalM    = flag.Bool("minimal",    false, "Use minimal decomposition")
var photoshopM  = flag.Bool("photoshop",  false, "Use photoshop luminosity method (default)")

var help = flag.Bool("help", false, "Display this help message")

func main() {
	flag.Parse()

	if *help {
		msg := "Usage: greyscale [opts]\n" +
			"\n" +
			"  Takes a png file from STDIN and prints a greyscale version to STDOUT.\n" +
			"  Can choose one of the many greyscaling methods.\n" +
			"\n" +
			"  --average           # Use average method\n" +
			"  --lightness         # Use lightness method\n" +
			"  --luminosity        # Use standard luminosity method\n" +
			"  --maximal           # Use maximal decompositon\n" +
			"  --minimal           # Use minimal decompositon\n" +
			"  --photoshop         # Use photoshop luminosity method (default)\n" +
			"  --help              # Display this help message\n" +
			"\n"

		fmt.Fprintf(os.Stderr, msg)
		os.Exit(0)
	}

	i := utils.ReadStdin()

	if *averageM {
		i = greyscale.Average(i)
	} else if *lightnessM {
    i = greyscale.Lightness(i)
	} else if *luminosityM {
		i = greyscale.Luminosity(i)
	} else if *maximalM {
		i = greyscale.Maximal(i)
	} else if *minimalM {
		i = greyscale.Minimal(i)
  } else {
		i = greyscale.Photoshop(i)
  }

	utils.WriteStdout(i)
}
