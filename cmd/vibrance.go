package cmd

import (
	"hawx.me/code/hadfield"
	"hawx.me/code/img/utils"
	"hawx.me/code/img/vibrance"
)

var (
	vibranceExp bool
	vibranceBy  float64
)

func Vibrance() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "vibrance [options]",
		Short: "adjust the vibrancy of an image",
		Long: `

	Vibrance adjusts the saturation of the least-saturated parts of an image to
	give the image more vibrancy.

    --exp        # Use inverse exponential method
    --by [n]     # Amount to adjust vibrancy by (default: 0.5)
`,
	}

	cmd.Run = runVibrance

	cmd.Flag.BoolVar(&vibranceExp, "exp", false, "")
	cmd.Flag.Float64Var(&vibranceBy, "by", 0.5, "")

	return cmd
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
