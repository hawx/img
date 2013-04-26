package main

import (
	"github.com/hawx/img/gamma"
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
)

var cmdGamma = &hadfield.Command{
	Usage: "gamma [options]",
	Short: "adjust image gamma",
Long: `
  Gamma takes an image from STDIN, adjusts the gamma and prints the result to
  STDOUT

    --auto         # Automatically alter gamma to "best" value (default)
    --by <n>       # Amount to adjust gamma by
    --undo         # Adjust by the reciprocal of the amount, instead
`,
}

var gammaAuto, gammaUndo bool
var gammaBy float64

func init() {
	cmdGamma.Run = runGamma

	cmdGamma.Flag.BoolVar(&gammaAuto,  "auto", false, "")
	cmdGamma.Flag.Float64Var(&gammaBy, "by",   1.8,   "")
	cmdGamma.Flag.BoolVar(&gammaUndo,  "undo", false, "")
}

func runGamma(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if gammaUndo { gammaBy = 1.0 / gammaBy }

	if utils.FlagVisited("by", cmd.Flag) {
		i = gamma.Adjust(i, gammaBy)
	} else {
		i = gamma.Auto(i)
	}

	utils.WriteStdout(i, data)
}
