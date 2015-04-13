package cmd

import (
	"github.com/hawx/hadfield"
	"github.com/hawx/img/greyscale"
	"github.com/hawx/img/utils"
)

var (
	greyscaleAverage, greyscaleLightness, greyscaleLuminosity bool
	greyscaleRed, greyscaleGreen, greyscaleBlue               bool
	greyscaleMaximal, greyscaleMinimal, greyscalePhotoshop    bool
)

func Greyscale() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "greyscale [options]",
		Short: "convert image to greyscale",
		Long: `
  Greyscale takes an image from STDIN, and prints to STDOUT a greyscale version

    --average        # Use average method
    --lightness      # Use lightness method
    --luminosity     # Use standard luminosity method
    --maximal        # Use maximal decomposition
    --minimal        # Use minimal decomposition
    --red            # Use the values of the red channel
    --green          # Use the values of the green channel
    --blue           # Use the values of the blue channel
`,
	}

	cmd.Run = runGreyscale

	cmd.Flag.BoolVar(&greyscaleAverage, "average", false, "")
	cmd.Flag.BoolVar(&greyscaleLightness, "lightness", false, "")
	cmd.Flag.BoolVar(&greyscaleLuminosity, "luminosity", false, "")
	cmd.Flag.BoolVar(&greyscaleMaximal, "maximal", false, "")
	cmd.Flag.BoolVar(&greyscaleMinimal, "minimal", false, "")
	cmd.Flag.BoolVar(&greyscaleRed, "red", false, "")
	cmd.Flag.BoolVar(&greyscaleGreen, "green", false, "")
	cmd.Flag.BoolVar(&greyscaleBlue, "blue", false, "")

	return cmd
}

func runGreyscale(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

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
	} else if greyscaleRed {
		i = greyscale.Red(i)
	} else if greyscaleGreen {
		i = greyscale.Green(i)
	} else if greyscaleBlue {
		i = greyscale.Blue(i)
	} else {
		i = greyscale.Greyscale(i)
	}

	utils.WriteStdout(i, data)
}
