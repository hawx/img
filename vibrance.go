package main

import (
	"github.com/hawx/hadfield"
	"github.com/hawx/img/utils"
	"github.com/hawx/img/vibrance"
)

var cmdVibrance = &hadfield.Command{
	Usage: "vibrance [options]",
	Short: "adjust the vibrancy of an image",
	Long: `

	Vibrance adjusts the saturation of the least-saturated parts of an image to
	give the image more vibrancy.

    --exp        # Use inverse exponential method
    --by [n]     # Amount to adjust vibrancy by (default: 0.5)
`,
}

var vibranceExp bool
var vibranceBy float64

func init() {
	cmdVibrance.Run = runVibrance

	cmdVibrance.Flag.BoolVar(&vibranceExp, "exp", false, "")
	cmdVibrance.Flag.Float64Var(&vibranceBy, "by", 0.5, "")
}

func runVibrance(cmd *hadfield.Command, args []string) {
	image, exif := utils.ReadStdin()

	if vibranceExp {
		image = vibrance.Exp(image, vibranceBy)
	} else {
		image = vibrance.Adjust(image, vibranceBy)
	}

	utils.WriteStdout(image, exif)
}
