package cmd

import (
	"hawx.me/code/hadfield"
	"hawx.me/code/img/sharpen"
	"hawx.me/code/img/utils"
)

var (
	sharpenRadius                                 int
	sharpenSigma, sharpenAmount, sharpenThreshold float64
	sharpenUnsharp                                bool
)

func Sharpen() *hadfield.Command {
	cmd := &hadfield.Command{
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

	cmd.Run = runSharpen

	// width=GetOptimalKernelWidth2D(radius,sigma);

	cmd.Flag.IntVar(&sharpenRadius, "radius", 1, "")
	cmd.Flag.Float64Var(&sharpenSigma, "sigma", 1.0, "")
	cmd.Flag.Float64Var(&sharpenAmount, "amount", 1.0, "")
	cmd.Flag.Float64Var(&sharpenThreshold, "threshold", 0.05, "")

	cmd.Flag.BoolVar(&sharpenUnsharp, "unsharp", false, "")

	return cmd
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
