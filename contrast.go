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

    --ratio <n>             # Ratio to shift contrast by (default: 6.0)
    --linear                # Use linear function
    --sigmoidal <midpoint>  # Use sigmoidal function (default: 0.5)
`,
}

var contrastLinear bool
var contrastRatio, contrastSigmoidal float64

func init() {
	cmdContrast.Run = runContrast

	cmdContrast.Flag.Float64Var(&contrastRatio, "ratio", 3.0, "")
	cmdContrast.Flag.BoolVar(&contrastLinear, "linear", false, "")
	cmdContrast.Flag.Float64Var(&contrastSigmoidal, "sigmoidal", 0.5, "")
}

func runContrast(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if contrastLinear {
		i  = contrast.Adjust(i, contrastRatio)
	} else {
		i = contrast.Sigmoidal(i, contrastRatio, contrastSigmoidal)
	}

	utils.WriteStdout(i)
}
