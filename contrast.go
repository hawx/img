package main

import (
	"github.com/hawx/img/contrast"
	"github.com/hawx/img/utils"
)

var cmdContrast = &Command{
	UsageLine: "contrast [options]",
	Short:     "adjust image contrast",
Long: `
  Contrast takes an image from STDIN, adjusts the contrast and prints the result
  to STDOUT

    --ratio <n>             # Ratio to shift contrast by (default: 3.0)
    --midpoint <n>          # Set where midtones fall in image (default: 0.5)

    --linear                # Use linear function (default)
    --sigmoidal             # Use sigmoidal function
`,
}

var contrastLinear, contrastSigmoidal bool
var contrastRatio, contrastMidpoint float64

func init() {
	cmdContrast.Run = runContrast

	cmdContrast.Flag.Float64Var(&contrastRatio,    "ratio",     3.0,   "")
	cmdContrast.Flag.Float64Var(&contrastMidpoint, "midpoint",  0.5,   "")

	cmdContrast.Flag.BoolVar(&contrastLinear,      "linear",    false, "")
	cmdContrast.Flag.BoolVar(&contrastSigmoidal,   "sigmoidal", false, "")
}

func runContrast(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if contrastSigmoidal {
		i = contrast.Sigmoidal(i, contrastRatio, contrastMidpoint)
	} else {
		i  = contrast.Adjust(i, contrastRatio)
	}

	utils.WriteStdout(i)
}
