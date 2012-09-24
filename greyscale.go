package main

import (
	"github.com/hawx/img/greyscale"
	"github.com/hawx/img/utils"
)

var cmdGreyscale = &Command{
	UsageLine: "greyscale [options]",
	Short:     "convert image to greyscale",
Long: `
  Greyscale takes an image from STDIN, and prints to STDOUT a greyscale version

    --average        # Use average method
    --lightness      # Use lightness method
    --luminosity     # Use standard luminosity method
    --maximal        # Use maximal decomposition
    --minimal        # Use minimal decomposition
    --photoshop      # Use photoshop luminosity method (default)
`,
}

var greyscaleAverage, greyscaleLightness, greyscaleLuminosity bool
var greyscaleMaximal, greyscaleMinimal, greyscalePhotoshop bool

func init() {
	cmdGreyscale.Run = runGreyscale

	cmdGreyscale.Flag.BoolVar(&greyscaleAverage,    "average",    false, "")
	cmdGreyscale.Flag.BoolVar(&greyscaleLightness,  "lightness",  false, "")
	cmdGreyscale.Flag.BoolVar(&greyscaleLuminosity, "luminosity", false, "")
	cmdGreyscale.Flag.BoolVar(&greyscaleMaximal,    "maximal",    false, "")
	cmdGreyscale.Flag.BoolVar(&greyscaleMinimal,    "minimal",    false, "")
	cmdGreyscale.Flag.BoolVar(&greyscalePhotoshop,  "photoshop",  false, "")
}

func runGreyscale(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if greyscaleAverage {
		i = greyscale.Average(i)
	} else if greyscaleLightness {
		i = greyscale.Lightness(i)
	} else if greyscaleLuminosity {
		i = greyscale.Luminosity(i)
	} else if greyscaleMaximal {
		i = greyscale.Maximal(i)
	} else if greyscaleMinimal {
		i = greyscale.Minimal(i)
	} else {
		i = greyscale.Photoshop(i)
	}

	utils.WriteStdout(i)
}
