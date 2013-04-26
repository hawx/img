package main

import (
	"github.com/hawx/img/contrast"
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
)

var cmdContrast = &hadfield.Command{
	Usage: "contrast [options]",
	Short: "adjust image contrast",
Long: `
  Contrast takes an image from STDIN, adjusts the contrast and prints the result
  to STDOUT

    --factor <n>            # Factor to shift contrast by (default: 1.0)
    --midpoint <n>          # Set where midtones fall in image (default: 0.5)

    --linear                # Use linear function
    --sigmoidal             # Use sigmoidal function
`,
}

var contrastLinear, contrastSigmoidal bool
var contrastFactor, contrastMidpoint float64

func init() {
	cmdContrast.Run = runContrast

	cmdContrast.Flag.Float64Var(&contrastFactor,   "factor",    1.0,   "")
	cmdContrast.Flag.Float64Var(&contrastMidpoint, "midpoint",  0.5,   "")

	cmdContrast.Flag.BoolVar(&contrastLinear,      "linear",    false, "")
	cmdContrast.Flag.BoolVar(&contrastSigmoidal,   "sigmoidal", false, "")
}

func runContrast(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if contrastSigmoidal {
		i = contrast.Sigmoidal(i, contrastFactor, contrastMidpoint)
	} else if contrastLinear {
		i  = contrast.Linear(i, contrastFactor)
	} else {
		i = contrast.Adjust(i, contrastFactor)
	}

	utils.WriteStdout(i, data)
}
