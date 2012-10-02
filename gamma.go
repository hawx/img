package main

import (
	"github.com/hawx/img/gamma"
	"github.com/hawx/img/utils"
)

var cmdGamma = &Command{
	UsageLine: "gamma [options]",
	Short:     "adjust image gamma",
Long: `
  Gamma takes an image from STDIN, adjusts the gamma and prints the result to
  STDOUT

    --auto         # Automatically alter gamma to "best" value (default)
    --by [n]       # Amount to adjust gamma by
`,
}

var gammaAuto bool
var gammaBy float64

func init() {
	cmdGamma.Run = runGamma

	cmdGamma.Flag.BoolVar(&gammaAuto, "auto", false, "")
	cmdGamma.Flag.Float64Var(&gammaBy, "by", 1.8, "")
}

func runGamma(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if gammaAuto {
		i = gamma.Auto(i)
	} else {
		i = gamma.Adjust(i, gammaBy)
	}
	utils.WriteStdout(i)
}
