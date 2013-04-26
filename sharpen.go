package main

import (
	"github.com/hawx/img/sharpen"
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
)

var cmdSharpen = &hadfield.Command{
	Usage: "sharpen [options]",
	Short: "sharpen an image",
Long: `
  Sharpen takes an image from STDIN, and prints a sharpened version to STDOUT.

    --radius <num>      # The radius of the Gaussian, not including center pixel
    --sigma <num>       # "Weightedness" of outer pixels of kernel
    --amount <num>      # Amount to sharpen by (between 0 and 1)
    --threshold <num>   # Fraction difference required to apply sharpen

    --unsharp
`,
}

var sharpenRadius int
var sharpenSigma, sharpenAmount, sharpenThreshold float64
var sharpenUnsharp bool

func init() {
	cmdSharpen.Run = runSharpen

	// width=GetOptimalKernelWidth2D(radius,sigma);

	cmdSharpen.Flag.IntVar(&sharpenRadius, "radius", 1, "")
	cmdSharpen.Flag.Float64Var(&sharpenSigma, "sigma", 1.0, "")
	cmdSharpen.Flag.Float64Var(&sharpenAmount, "amount", 1.0, "")
	cmdSharpen.Flag.Float64Var(&sharpenThreshold, "threshold", 0.05, "")

	cmdSharpen.Flag.BoolVar(&sharpenUnsharp, "unsharp", false, "")
}

func runSharpen(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if sharpenUnsharp {
		i = sharpen.UnsharpMask(i, sharpenRadius, sharpenSigma, sharpenAmount, sharpenThreshold)
	} else {
		i = sharpen.Sharpen(i, sharpenRadius, sharpenSigma)
	}

	utils.WriteStdout(i, data)
}
